package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

func makeRequest(client *http.Client, baseURL string, key string, value string, wg *sync.WaitGroup, totalTime *time.Duration) {
	defer wg.Done()

	url := baseURL + key
	req, err := http.NewRequest("PUT", url, bytes.NewBufferString(value))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Content-Type", "text/plain")

	startTime := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	timeToFirstByte := time.Since(startTime)
	// fmt.Printf("Request for key %s - Time taken for the first byte: %v\n", key, timeToFirstByte)

	// Update totalTime with the time taken for this request
	*totalTime += timeToFirstByte
}

func main() {
	baseURL := "http://128.110.216.12:8060/"
	numRequests := 10000
	totalTime := time.Duration(0)
	client := &http.Client{}

	var wg sync.WaitGroup
	wg.Add(numRequests)

	for i := 1; i <= numRequests; i++ {
		key := fmt.Sprintf("my-key-%d", i)
		value := fmt.Sprintf("my-value-%d", i)
		go makeRequest(client, baseURL, key, value, &wg, &totalTime)
	}

	wg.Wait()
	fmt.Println("All requests completed.")

	averageTime := totalTime / time.Duration(numRequests)
	fmt.Println("Average Time to First Byte:", averageTime)
}
