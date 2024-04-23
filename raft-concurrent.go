package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"sync"
	"os"
	"strconv"
	"strings"
)


type Barrier struct {
	count  int
	n      int
	mu     sync.Mutex
	cond   *sync.Cond
}

func NewBarrier(n int) *Barrier {
	b := &Barrier{
		count: 0,
		n:     n,
	}
	b.cond = sync.NewCond(&b.mu)
	return b
}

func (b *Barrier) Wait() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.count++
	if b.count == b.n {
		b.count = 0
		b.cond.Broadcast()
	} else {
		b.cond.Wait()
	}
}


func main() {

	numRequests := 0
	if len(os.Args) < 2 {
		numRequests = 5
		fmt.Println("Usage: go run main.go <integer>")
	} else {
		arg := os.Args[1]
		cmdNum, err := strconv.Atoi(arg)
		if err != nil {
			numRequests = 5
			fmt.Println("Error: conversion", err)
			
		} else {
			numRequests = cmdNum
		}
		
	}

	// Base URL for the PUT request
	baseURL := "http://128.110.216.12:8050/"


	// Initialize total time
	totalTime := time.Duration(0)
	maxTime := time.Duration(0)

		// Number of requests to make

	



	durarr := make([]time.Duration, numRequests)
	// Create a new HTTP client
	client := &http.Client{}

	barrier := NewBarrier(numRequests)
	var wg sync.WaitGroup

	// Loop for making requests
	for i := 0; i < numRequests; i++ {

		wg.Add(1)
		go func (id int) {
			
			barrier.Wait()
			defer wg.Done()

			
			// Generate unique key and value
			key := fmt.Sprintf("my-key-%d", id)
			value := fmt.Sprintf("my-value-%d", id)
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

			done := false
			for !done {
				// Perform the PUT request
				resp, err := client.Do(req)
				if err != nil {
					fmt.Println("Error making request:", err)
					continue
				}
				defer resp.Body.Close()

				// Read response body
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					fmt.Println("Error reading response body:", err)
					return
				}
				if strings.TrimSpace((string(body))) != "done" {
					// fmt.Println("Response not done:", string(body))

				} else {
					done = true
				}

			}
			


			// Calculate the time taken for the first byte
			timeToFirstByte := time.Since(startTime)

			durarr[id] = timeToFirstByte
		}(i)

		
		// Aggregate the time taken

		// Print response status code and time taken for the first byte
		// fmt.Printf("Request %d - Time taken for the first byte: %v\n", i, timeToFirstByte)
	}

	wg.Wait()
	for i := 0; i < numRequests; i++{
		if durarr[i] > maxTime {
			maxTime = durarr[i]
		}

		totalTime += durarr[i]
	}

	// Calculate the average time to first byte
	averageTime := totalTime / time.Duration(numRequests)

	// Print the average time to first byte
	fmt.Println("Average Time to First Byte:", averageTime)
	fmt.Println("Max Time:", maxTime)

}