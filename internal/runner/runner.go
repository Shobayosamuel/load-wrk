package runner

import (
    "net/http"
    "sync"
    "time"
    "strings"
    "github.com/Shobayosamuel/load-wrk/internal/metrics"
)

type Result struct {
    StatusCode int
    Duration   time.Duration
    Error      error
}

func StartWorkers(targetURL string, method string, totalRequests int, concurrency int, rate int, body string, headers []string) []Result {
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
                req, _ := http.NewRequest(method, targetURL, strings.NewReader(body))
                for _, h := range headers {
                    parts := strings.SplitN(h, ":", 2)
                    if len(parts) == 2 {
                        req.Header.Set(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
                    }
                }

                resp, err := client.Do(req)
                duration := time.Since(start)
				metrics.RequestsTotal.Inc()
                if err != nil {
                    resultChan <- Result{StatusCode: 0, Duration: duration, Error: err}
                    continue
                }
				metrics.RequestsSuccess.Inc()
				metrics.LatencyHistogram.Observe(duration.Seconds())
                resp.Body.Close()
                resultChan <- Result{StatusCode: resp.StatusCode, Duration: duration}
            }
        }(i)
    }

    // Send jobs
    if rate > 0 {
        ticker := time.NewTicker(time.Second / time.Duration(rate))
        for i := 0; i < totalRequests; i++ {
            <-ticker.C
            jobs <- i
        }
        ticker.Stop()
    } else {
        for i := 0; i < totalRequests; i++ {
            jobs <- i
        }
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
