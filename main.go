package main

import (
	"CheckGraphite"
	"GraphiteData"
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
)

var opts struct {
	Host       string  `short:"H" long:"host" description:"Graphite host url (required)" value-name:"HOST" required:"true"`
	Metric     string  `short:"m" long:"metric" description:"Graphite metric name (required)" value-name:"METRIC" required:"true"`
	NullIsZero bool    `short:"z" long:"zero" description:"Convert 'None' values to 0."`
	Duration   int     `short:"d" long:"duration" description:"Number of minutes of data to aggregate." value-name:"SECONDS" default:"10"`
	Function   string  `short:"f" long:"function" description:"The aggregation function to apply." value-name:"(min|max|avg|sum|last)" default:"last"`
	Warning    float64 `short:"w" long:"warning" description:"Warning threshold of aggregated value." value-name:"WARNING"`
	Critical   float64 `short:"c" long:"critical" description:"Critical threshold of aggregated value." value-name:"CRITICAL"`
	Invert     bool    `short:"i" long:"invert" description:"Invert thresholds to alert below metric value."`
}

func main() {

	// Get Options
	parser := flags.NewParser(&opts, flags.HelpFlag)
	if _, err := parser.Parse(); err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(3)
	}
	if (opts.Warning > opts.Critical) && opts.Invert == false {
		fmt.Printf("UNKNOWN: Warning threshold must be less than Critical threshold.")
		os.Exit(3)
	}
	if (opts.Warning < opts.Critical) && opts.Invert == true {
		fmt.Printf("UNKNOWN: Warning threshold must be greater than Critical threshold.")
		os.Exit(3)
	}

	// Get Metrics from Graphite
	graphiteData, err := GraphiteData.GetMetrics(opts.Host, opts.Metric, opts.Duration)
	if err != nil {
		fmt.Printf("UNKNOWN: %s\n", err)
		os.Exit(3)
	}
	if len(graphiteData) > 1 {
		fmt.Printf("UNKNOWN: Graphite query must return a single Target.\n")
		os.Exit(3)
	}

	// Aggregate Metrics
	metrics, err := CheckGraphite.New(graphiteData[0], opts.NullIsZero, opts.Function)
	if err != nil {
		fmt.Printf("UNKNOWN: %s\n", err)
		os.Exit(3)
	}

	// Do Alerts
	msg, exit := metrics.DoAlerts(opts.Warning, opts.Critical, opts.Invert)

	// Print and Exit
	fmt.Printf("%s\n", msg)
	os.Exit(exit)

}
