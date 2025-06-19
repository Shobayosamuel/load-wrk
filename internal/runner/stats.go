package runner

import (
    "sort"
    "time"
)

type Stats struct {
    Total        int
    Success      int
    Failed       int
    AvgLatency   time.Duration
    P95Latency   time.Duration
    MaxLatency   time.Duration
    RequestsPerSec float64
}

func CalculateStats(results []Result) Stats {
    var (
        durations   []time.Duration
        total       = len(results)
        success     int
        maxDuration time.Duration
        totalDuration time.Duration
    )
	// Get the duration results
    for _, r := range results {
        if r.Error == nil {
            success++
        }
        durations = append(durations, r.Duration)
        totalDuration += r.Duration
        if r.Duration > maxDuration {
            maxDuration = r.Duration
        }
    }

    sort.Slice(durations, func(i, j int) bool {
        return durations[i] < durations[j]
    })
	// Get the 95th percentile
    p95Index := int(float64(len(durations)) * 0.95)
    if p95Index >= len(durations) {
        p95Index = len(durations) - 1
    }

    avg := time.Duration(0)
    if total > 0 {
        avg = totalDuration / time.Duration(total)
    }

    rps := float64(total) / totalDuration.Seconds()

    return Stats{
        Total:        total,
        Success:      success,
        Failed:       total - success,
        AvgLatency:   avg,
        P95Latency:   durations[p95Index],
        MaxLatency:   maxDuration,
        RequestsPerSec: rps,
    }
}
