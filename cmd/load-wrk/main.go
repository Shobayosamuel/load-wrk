package main

import (
	"fmt"
	"os"

	"github.com/Shobayosamuel/load-wrk/internal/metrics"
	"github.com/Shobayosamuel/load-wrk/internal/output"
	"github.com/Shobayosamuel/load-wrk/internal/runner"
	"github.com/spf13/cobra"
)
func main() {
    var (
        concurrency int
        requests    int
        method      string
        prometheus  bool
        rate        int
        body string
        headers []string
    )

    var rootCmd = &cobra.Command{
        Use:   "load-wrk",
        Short: "A simple Go load testing tool",
    }
    var runCmd = &cobra.Command{
        Use:   "run [url]",
        Short: "Run a load test against a target URL",
        Args:  cobra.ExactArgs(1),
        Run: func(cmd *cobra.Command, args []string) {
            url := args[0]
            if prometheus {
                metrics.InitMetrics()
                fmt.Println("Starting Prometheus metrics server at :2112/metrics")
                metrics.StartPrometheusServer()
            }

            results := runner.StartWorkers(url, method, requests, concurrency, rate, body, headers)
            stats := runner.CalculateStats(results)
            output.PrintSummary(stats)
			output.PrintLatencyGraph(stats.Durations)
            if prometheus {
                fmt.Println("Server running. Press Ctrl+C to exit.")
                select {} // block forever
            }
        },
    }


    runCmd.Flags().StringVar(&body, "body", "", "Request body for POST/PUT requests")
    runCmd.Flags().StringArrayVar(&headers, "header", []string{}, "Custom headers (key:value format)")

    runCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 10, "Number of concurrent workers")
    runCmd.Flags().IntVarP(&requests, "requests", "n", 100, "Total number of requests")
    runCmd.Flags().StringVarP(&method, "method", "m", "GET", "HTTP method to use")
	runCmd.Flags().BoolVar(&prometheus, "prometheus", false, "Expose Prometheus metrics at :2112/metrics")
    runCmd.Flags().IntVarP(&rate, "rate", "r", 0, "Rate limit (requests per second, 0 = unlimited)")

    rootCmd.AddCommand(runCmd)

    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}