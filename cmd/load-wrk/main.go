package main

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
    "github.com/Shobayosamuel/load-wrk/internal/runner"
    "github.com/Shobayosamuel/load-wrk/internal/output"
)

func main() {
    var concurrency, requests int
    var method string

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

            results := runner.StartWorkers(url, method, requests, concurrency)
            stats := runner.CalculateStats(results)
            output.PrintSummary(stats)
        },
    }

    runCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 10, "Number of concurrent workers")
    runCmd.Flags().IntVarP(&requests, "requests", "n", 100, "Total number of requests")
    runCmd.Flags().StringVarP(&method, "method", "m", "GET", "HTTP method to use")

    rootCmd.AddCommand(runCmd)

    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}