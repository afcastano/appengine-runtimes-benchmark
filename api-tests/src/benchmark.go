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
	Runtime  string
	TestType string
	Request  string
	Decode   decodeJsonFunc
}

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

func loadKeys(testType string) ([]string, error) {
	if testType == "query" {
		return makeRange(400000), nil
	} else if testType == "key" {
		return loadEntityNames()
	} else {
		return nil, errors.New(fmt.Sprintf("Test type %s is not defined", testType))
	}
}

func loadEntityNames() ([]string, error) {
	start := time.Now()
	resp, err := http.Get("https://gae-benchmark.appspot.com/names")

	if err != nil {
		return nil, err
	}

	var keys []string
	err = json.NewDecoder(resp.Body).Decode(&keys)

	if err != nil {
		return nil, err
	}

	fmt.Printf("%.2fs fetch %d keys \n", time.Since(start).Seconds(), len(keys))
	return keys, nil
}

func getEnvironmentContext(runtime string, testType string) (EnvironmentContext, error) {
	var context EnvironmentContext
	var requestStr string
	var decoder decodeJsonFunc = defaultJsonDecoder

	if runtime == "spring" {
		if testType == "key" {
			requestStr = "https://java-spring-dot-gae-benchmark.appspot.com/entity/%s"
		} else if testType == "query" {
			requestStr = "https://java-spring-dot-gae-benchmark.appspot.com/entities/%s"
		} else {
			return context, errors.New(fmt.Sprintf("Test type %s does not exist for runtime %s", testType, runtime))
		}

	} else if runtime == "thundr" {
		if testType == "key" {
			requestStr = "https://java-thundr-dot-gae-benchmark.appspot.com/entity/%s"
		} else if testType == "query" {
			requestStr = "https://java-thundr-dot-gae-benchmark.appspot.com/entity/greaterThan/%s"
		} else {
			return context, errors.New(fmt.Sprintf("Test type %s does not exist for runtime %s", testType, runtime))
		}

	} else if runtime == "nest" {
		if testType == "key" {
			requestStr = "https://nodejs-dot-gae-benchmark.appspot.com/api/graphql?query={getDummyById(id:\"%s\"){id,random1,random2}}"
		} else if testType == "query" {
			requestStr = "https://nodejs-dot-gae-benchmark.appspot.com/api/graphql?query={dummies(index:%s){id,random1,random2}}"
			decoder = nestJsonDecoder
		} else {
			return context, errors.New(fmt.Sprintf("Test type %s does not exist for runtime %s", testType, runtime))
		}
	} else if runtime == "go" {
		if testType == "query" {
			requestStr = "https://go-dot-gae-benchmark.appspot.com/entities?random2=%s"
		} else {
			return context, errors.New(fmt.Sprintf("Test type %s does not exist for runtime %s", testType, runtime))
		}

	} else if runtime == "express" {
		if testType == "query" {
			requestStr = "https://node-express-dot-gae-benchmark.appspot.com/entities/%s"
		} else {
			return context, errors.New(fmt.Sprintf("Test type %s does not exist for runtime %s", testType, runtime))
		}

	} else {
		return context, errors.New(fmt.Sprintf("Runtime %s not configured", runtime))
	}

	fmt.Printf("Processing qurery %s \n", requestStr)
	return EnvironmentContext{Runtime: runtime, Decode: decoder, TestType: testType, Request: requestStr}, nil
}

func loopRequest(context EnvironmentContext, keys []string, chanId int, repeatCount int, ch chan<- ChannelResponse) {
	fmt.Printf("Starting goroutine %d\n", chanId)
	start := time.Now()
	successCount := 0
	errCount := 0
	totalCount := 0

	for i := 0; i < repeatCount; i++ {
		succ, err := loopAllValues(chanId, keys, context)
		successCount += succ
		errCount += err
		totalCount += len(keys)
		fmt.Printf("Channel %d has finished loop %d\n", chanId, i)
	}

	secs := time.Since(start).Seconds()
	response := ChannelResponse{ChanId: chanId, TotalRequests: totalCount, TotalErrors: errCount, TimeTaken: secs}
	ch <- response
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

func loopAllValues(chanId int, keys []string, context EnvironmentContext) (int, int) {
	errCount := 0
	successCount := 0

	for idx, key := range keys {
		resp, err := http.Get(fmt.Sprintf(context.Request, key))

		if idx%100 == 0 {
			fmt.Printf("Channel %d has processed %d requests\n", chanId, idx)
		}

		if err != nil {
			errCount += 1
			continue
		}

		entities, err := context.Decode(resp.Body)

		if err != nil {
			fmt.Printf("WARNING Channel %d had an error while decoding values in %s\n", chanId, key)
			errCount += 1
			continue
		}

		if len(entities) == 0 {
			fmt.Printf("WARNING Channel %d didn't load any entity for key %s\n", chanId, key)
		}

		successCount += 1
	}

	return successCount, errCount
}

func main() {
	args := os.Args[1:]
	runtime := args[0]
	testType := args[1]
	routines, _ := strconv.Atoi(args[2])

	fmt.Printf("Benchmarking runtime %s type %s with %d routines\n", runtime, testType, routines)

	start := time.Now()

	context, err := getEnvironmentContext(runtime, testType)

	if err != nil {
		panic(err)
	}

	keys, err := loadKeys(testType)
	if err != nil {
		panic(err)
	}

	ch := make(chan ChannelResponse)
	numberOfLoops := 10
	// Start goroutines fetching keys one by one in a loop 10 times.
	for channelId := 0; channelId < routines; channelId++ {
		go loopRequest(context, keys, channelId, numberOfLoops, ch)
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
