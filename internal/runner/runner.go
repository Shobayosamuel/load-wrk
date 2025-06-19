package runner

import (
    "net/http"
    "sync"
    "time"
)

type Result struct {
    StatusCode int
    Duration   time.Duration
    Error      error
}

func StartWorkers(targetURL string, method string, totalRequests int, concurrency int) []Result {
    jobs := make(chan int, totalRequests)
    results := make([]Result, 0, totalRequests)
    resultChan := make(chan Result, totalRequests)

    var wg sync.WaitGroup

    // Spawn N worker goroutines
    for i := 0; i < concurrency; i++ {
        wg.Add(1)
        go func(workerID int) {
            defer wg.Done()
            client := &http.Client{}

            for range jobs {
                start := time.Now()
                req, _ := http.NewRequest(method, targetURL, nil)

                resp, err := client.Do(req)
                duration := time.Since(start)

                if err != nil {
                    resultChan <- Result{StatusCode: 0, Duration: duration, Error: err}
                    continue
                }
                resp.Body.Close()
                resultChan <- Result{StatusCode: resp.StatusCode, Duration: duration}
            }
        }(i)
    }

    // Send jobs
    for i := 0; i < totalRequests; i++ {
        jobs <- i
    }
    close(jobs)

    // Wait for all workers to finish
    go func() {
        wg.Wait()
        close(resultChan)
    }()

    // Collect results
    for r := range resultChan {
        results = append(results, r)
    }

    return results
}
