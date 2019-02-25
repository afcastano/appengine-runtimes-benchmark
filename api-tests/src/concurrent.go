package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func loadKeys() ([]string, error) {
	start := time.Now()
	resp, err := http.Get("https://runtimes-benchmark.appspot.com/load")

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

func requestKeysOneByOne(chanId int, keys []string, ch chan<- string) {
	fmt.Printf("Starting goroutine %d\n", chanId)
	start := time.Now()
	errCount := 0
	successCount := 0

	for idx, key := range keys {
		// resp, err := http.Get(fmt.Sprintf("https://java-spring-dot-runtimes-benchmark.appspot.com/entity/%s", key))
		resp, err := http.Get(fmt.Sprintf("https://java-thundr-dot-runtimes-benchmark.appspot.com/entity/%s", key))
		if idx%25 == 0 {
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

	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("Channel %d finished in %.2f with %d success calls and %d errors\n", chanId, secs, successCount, errCount)
}

func main() {
	start := time.Now()

	keys, err := loadKeys()

	if err != nil {
		fmt.Println("Error $s", err)
		return
	}

	ch := make(chan string)
	// Start 20 goroutines fetching keys one by one in a loop 3000 times.
	for i := 0; i < 100; i++ {
		go requestKeysOneByOne(i, keys, ch)
	}

	for i := 0; i < 100; i++ {
		fmt.Println(<-ch)
	}

	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}
