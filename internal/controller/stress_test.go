// +build stress

package controller

import (
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const (
	numRequests = 1000
	numWorkers  = 10
	url         = "http://127.0.0.1:3030/v1/hello-world" // replace with url of your choice
)

// make sure the server is running first before running this test else it will throw an error.
//  Use 'make run' to start the server
func TestServer(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	var (
		wg            sync.WaitGroup
		totalDuration int64
		successCount  int64
		errorCount    int64
	)

	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
		},
	}
	startTest := time.Now()

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(i, &wg, client, &totalDuration, &successCount, &errorCount)
	}

	wg.Wait()
	endTest := time.Since(startTest)

	fmt.Printf("Total time for test(seconds): %v\n", endTest.Seconds())
	fmt.Printf("Average response time(seconds): %v\n", time.Duration(totalDuration/numRequests).Seconds())
	fmt.Printf("Successful requests: %d\n", successCount)
	fmt.Printf("Failed requests: %d\n", errorCount)
	fmt.Printf("Requests per second: %.2f\n", float64(numRequests)/endTest.Seconds())
}

func worker(id int, wg *sync.WaitGroup, client *http.Client, totalDuration, successCount, errorCount *int64) {
	defer wg.Done()

	for i := 0; i < numRequests/numWorkers; i++ {
		start := time.Now()
		resp, err := client.Get(url)
		if err != nil {
			fmt.Printf("Worker %d - Error: %v\n", id, err)
			atomic.AddInt64(errorCount, 1)
			continue
		}
		resp.Body.Close()
		duration := time.Since(start)
		atomic.AddInt64(totalDuration, int64(duration))
		atomic.AddInt64(successCount, 1)
	}
}
