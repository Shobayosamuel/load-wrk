package output

import (
    "fmt"
    "os"

    "github.com/olekukonko/tablewriter"
    "github.com/Shobayosamuel/load-wrk/internal/runner"
)

func PrintSummary(stats runner.Stats) {
    table := tablewriter.NewWriter(os.Stdout)
    table.SetHeader([]string{"Metric", "Value"})

    table.Append([]string{"Total Requests", fmt.Sprintf("%d", stats.Total)})
    table.Append([]string{"Success", fmt.Sprintf("%d", stats.Success)})
    table.Append([]string{"Failed", fmt.Sprintf("%d", stats.Failed)})
    table.Append([]string{"Avg Latency", stats.AvgLatency.String()})
    table.Append([]string{"P95 Latency", stats.P95Latency.String()})
    table.Append([]string{"Max Latency", stats.MaxLatency.String()})
    table.Append([]string{"Requests/sec", fmt.Sprintf("%.2f", stats.RequestsPerSec)})

    table.Render()
}
