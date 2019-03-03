package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

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
	resp, err := http.Get("https://gae-runtimes-benchmark.appspot.com/names")

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

func getRequest(runtime string, testType string) (string, error) {
	var requestStr string
	if runtime == "spring" {
		if testType == "key" {
			requestStr = "https://java-spring-dot-gae-runtimes-benchmark.appspot.com/entity/%s"
		} else if testType == "query" {
			requestStr = "https://java-spring-dot-gae-runtimes-benchmark.appspot.com/entities/%s"
		} else {
			return "", errors.New(fmt.Sprintf("Test type %s does not exist for runtime %s", testType, runtime))
		}

	} else if runtime == "thundr" {
		if testType == "key" {
			requestStr = "https://java-thundr-dot-gae-runtimes-benchmark.appspot.com/entity/%s"
		} else if testType == "query" {
			requestStr = "https://java-thundr-dot-gae-runtimes-benchmark.appspot.com/entity/greaterThan/%s"
			fmt.Printf("Processing qurery %s \n", requestStr)
		} else {
			return "", errors.New(fmt.Sprintf("Test type %s does not exist for runtime %s", testType, runtime))
		}

	} else if runtime == "node" {
		if testType == "key" {
			requestStr = "https://nodejs-dot-gae-runtimes-benchmark.appspot.com/api/graphql?query={getDummyById(id:\"%s\"){id,random1,random2}}"
		} else if testType == "query" {
			requestStr = "https://nodejs-dot-gae-runtimes-benchmark.appspot.com/api/graphql?query={dummies(index:%s){id,random1,random2}}"
		} else {
			return "", errors.New(fmt.Sprintf("Test type %s does not exist for runtime %s", testType, runtime))
		}
	} else {
		return "", errors.New(fmt.Sprintf("Runtime %s not configured", runtime))
	}

	return requestStr, nil
}

func loopRequest(request string, keys []string, chanId int, repeatCount int, ch chan<- string) {
	fmt.Printf("Starting goroutine %d\n", chanId)
	start := time.Now()
	successCount := 0
	errCount := 0

	for i := 0; i < repeatCount; i++ {
		succ, err := loopAllValues(chanId, keys, request)
		successCount += succ
		errCount += err
		fmt.Printf("Channel %d has finished loop %d\n", chanId, i)
	}

	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("Channel %d finished in %.2f with %d success calls and %d errors\n", chanId, secs, successCount, errCount)
}

func loopAllValues(chanId int, keys []string, request string) (int, int) {
	errCount := 0
	successCount := 0

	for idx, key := range keys {
		resp, err := http.Get(fmt.Sprintf(request, key))

		if idx%100 == 0 {
			fmt.Printf("Channel %d has processed %d requests\n", chanId, idx)
		}

		if err != nil {
			errCount += 1
			continue
		}

		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			errCount += 1
			continue
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

	request, err := getRequest(runtime, testType)

	if err != nil {
		panic(err)
	}

	keys, err := loadKeys(testType)
	if err != nil {
		panic(err)
	}

	ch := make(chan string)
	// Start goroutines fetching keys one by one in a loop 10 times.
	for i := 0; i < routines; i++ {
		go loopRequest(request, keys, i, 5, ch)
	}

	for i := 0; i < routines; i++ {
		fmt.Println(<-ch)
	}

	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}
