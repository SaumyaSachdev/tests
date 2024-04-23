package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

func makeRequests(client *http.Client, baseURL string, startKey int, endKey int, wg *sync.WaitGroup, totalTime *time.Duration) {
	defer wg.Done()
	incorrectRes := 0

	for i := startKey; i <= endKey; i++ {
		key := fmt.Sprintf("my-key-%d", i)
		value := fmt.Sprintf("my-value-%d", i)

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

				// Check response content and content-type
		if resp.StatusCode == http.StatusOK {
			contentType := resp.Header.Get("Content-Type")
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error reading response body:", err)
				return
			}
			fmt.Println("Response body: ", string(body));
			if string(body) != "done" || contentType != "text/plain" {
				// fmt.Println("Invalid response content or content-type.")
				
				incorrectRes++
				fmt.Printf("incorrect responses number i: %d total: %d\n", i, incorrectRes);
			}
		} else {
			fmt.Printf("Error response code: %d\n", resp.StatusCode)
			return
		}
		// Update totalTime with the time taken for this request
		
		*totalTime += timeToFirstByte
	}
}

func main() {
	baseURL := "http://128.110.216.12:8060/"
	numRequests := 100
	numCores := 2
	totalTime := time.Duration(0)
	client := &http.Client{}

	var wg sync.WaitGroup
	wg.Add(numCores)

	requestsPerCore := numRequests / numCores

	for core := 0; core < numCores; core++ {
		startKey := core*requestsPerCore + 1
		endKey := startKey + requestsPerCore - 1
		if core == numCores-1 {
			endKey = numRequests
		}
		go makeRequests(client, baseURL, startKey, endKey, &wg, &totalTime)
	}

	wg.Wait()
	fmt.Println("All requests completed.")

	averageTime := totalTime / time.Duration(numRequests)
	fmt.Println("Average Time to First Byte:", averageTime)
}
