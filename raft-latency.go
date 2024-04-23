package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	// Base URL for the PUT request
	baseURL := "http://128.110.216.12:8050/"

	// Number of requests to make
	numRequests := 10000

	// Initialize total time
	totalTime := time.Duration(0)

	// Create a new HTTP client
	client := &http.Client{}

	// Loop for making requests
	for i := 1; i <= numRequests; i++ {
		// Generate unique key and value
		key := fmt.Sprintf("my-key-%d", i)
		value := fmt.Sprintf("my-value-%d", i)

		// URL for the PUT request with unique key
		url := baseURL + key

		// Create a new request with PUT method and request body
		req, err := http.NewRequest("PUT", url, bytes.NewBufferString(value))
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}

		// Set Content-Type header
		req.Header.Set("Content-Type", "text/plain")

		// Record the start time
		startTime := time.Now()

		// Perform the PUT request
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error making request:", err)
			return
		}
		defer resp.Body.Close()

		// Read response body
		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return
		}

		// Calculate the time taken for the first byte
		timeToFirstByte := time.Since(startTime)

		// Aggregate the time taken
		totalTime += timeToFirstByte

		// Print response status code and time taken for the first byte
		// fmt.Printf("Request %d - Time taken for the first byte: %v\n", i, timeToFirstByte)
	}

	// Calculate the average time to first byte
	averageTime := totalTime / time.Duration(numRequests)

	// Print the average time to first byte
	fmt.Println("Average Time to First Byte:", averageTime)
}

