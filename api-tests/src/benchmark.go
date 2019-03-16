package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type GQLResult struct {
	Data GQLDummies `json:"data"`
}

type GQLDummies struct {
	Dummies []DummyEntity `json:"dummies"`
}

type DummyEntity struct {
	Id      string `json:"id"`
	Random1 string `json:"random1"`
	Random2 int    `json:"random2"`
}

type ChannelResponse struct {
	ChanId        int
	TotalRequests int
	TotalErrors   int
	TimeTaken     float64
}

type decodeJsonFunc func(body io.ReadCloser) ([]DummyEntity, error)

type EnvironmentContext struct {
	Runtime string
	Request string
	Decode  decodeJsonFunc
}

const MAX_FILTER_IDX = 400000
const LOOP_STEP = 1000

func makeRange(max int) []string {
	step := 1000
	arraySize := max / step
	a := make([]string, arraySize)

	for i := 0; i < arraySize; i++ {
		a[i] = strconv.Itoa(i * step)
	}

	fmt.Printf("Generated %d integers \n", len(a))
	return a
}

func defaultJsonDecoder(body io.ReadCloser) ([]DummyEntity, error) {
	var entities []DummyEntity
	err := json.NewDecoder(body).Decode(&entities)
	if err != nil {
		return nil, err
	}

	return entities, nil
}

func nestJsonDecoder(body io.ReadCloser) ([]DummyEntity, error) {
	var result GQLResult
	err := json.NewDecoder(body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result.Data.Dummies, nil
}

func getEnvironmentContext(runtime string) (EnvironmentContext, error) {
	var context EnvironmentContext
	var requestStr string
	var decoder decodeJsonFunc = defaultJsonDecoder

	if runtime == "spring" {
		requestStr = "https://java-spring-dot-gae-benchmark.appspot.com/entities/%d"
	} else if runtime == "thundr" {
		requestStr = "https://java-thundr-dot-gae-benchmark.appspot.com/entity/greaterThan/%d"
	} else if runtime == "nest" {
		requestStr = "https://nodejs-dot-gae-benchmark.appspot.com/api/graphql?query={dummies(index:%d){id,random1,random2}}"
		decoder = nestJsonDecoder
	} else if runtime == "go" {
		requestStr = "https://go-dot-gae-benchmark.appspot.com/entities?random2=%d"
	} else if runtime == "express" {
		requestStr = "https://node-express-dot-gae-benchmark.appspot.com/entities/%d"
	} else {
		return context, errors.New(fmt.Sprintf("Runtime %s not configured", runtime))
	}

	fmt.Printf("Processing qurery %s \n", requestStr)
	return EnvironmentContext{Runtime: runtime, Decode: decoder, Request: requestStr}, nil
}

func startChannel(context EnvironmentContext, chanId int, repeatCount int, ch chan<- ChannelResponse) {
	fmt.Printf("Starting goroutine %d\n", chanId)
	start := time.Now()
	successCount := 0
	errCount := 0
	totalCount := 0

	for i := 0; i < repeatCount; i++ {
		succ, err := makeRequests(chanId, context)
		successCount += succ
		errCount += err
		totalCount += MAX_FILTER_IDX / LOOP_STEP
		fmt.Printf("Channel %d has finished loop %d\n", chanId, i)
	}

	secs := time.Since(start).Seconds()
	response := ChannelResponse{ChanId: chanId, TotalRequests: totalCount, TotalErrors: errCount, TimeTaken: secs}
	ch <- response
}

func makeRequests(chanId int, context EnvironmentContext) (int, int) {
	errCount := 0
	successCount := 0
	for idx := 0; idx < MAX_FILTER_IDX; idx += LOOP_STEP {
		resp, err := http.Get(fmt.Sprintf(context.Request, idx))
		reqNo := idx / LOOP_STEP
		if (reqNo % 100) == 0 {
			fmt.Printf("Channel %d has processed %d requests\n", chanId, reqNo)
		}

		if err != nil {
			errCount += 1
			continue
		}

		entities, err := context.Decode(resp.Body)

		if err != nil {
			fmt.Printf("WARNING Channel %d had an error while decoding values in %d\n", chanId, idx)
			errCount += 1
			continue
		}

		if len(entities) == 0 {
			fmt.Printf("WARNING Channel %d didn't load any entity for key %d\n", chanId, idx)
		}

		successCount += 1
	}
	return successCount, errCount
}

func main() {
	args := os.Args[1:]
	runtime := args[0]
	routines, _ := strconv.Atoi(args[1])

	fmt.Printf("Benchmarking runtime %s with %d routines\n", runtime, routines)

	start := time.Now()

	context, err := getEnvironmentContext(runtime)

	if err != nil {
		panic(err)
	}

	ch := make(chan ChannelResponse)
	numberOfLoops := 10
	// Start goroutines fetching keys one by one in a loop 10 times.
	for channelId := 0; channelId < routines; channelId++ {
		go startChannel(context, channelId, numberOfLoops, ch)
	}

	totalRequests := 0
	for i := 0; i < routines; i++ {
		response := <-ch
		totalRequests += response.TotalRequests
		fmt.Printf("Channel %d finished on %.2f req/sec average\n", response.ChanId, float64(response.TotalRequests)/response.TimeTaken)
	}

	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
	fmt.Printf("Average request per second: %.2f\n", float64(totalRequests)/time.Since(start).Seconds())
}
