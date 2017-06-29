package CheckGraphite

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/pablojudd/go-graphite-getmetrics"
)

var graphiteData = []GraphiteData.GraphiteMetrics{
	{
		Target: "collectd.graphite.load.load.longterm",
		Datapoints: [][2]GraphiteData.NullFloat64{
			{
				GraphiteData.NullFloat64{sql.NullFloat64{Float64: 0, Valid: false}},
				GraphiteData.NullFloat64{sql.NullFloat64{Float64: 1.49660502e+09, Valid: true}},
			},
			{
				GraphiteData.NullFloat64{sql.NullFloat64{Float64: 0, Valid: false}},
				GraphiteData.NullFloat64{sql.NullFloat64{Float64: 1.49660508e+09, Valid: true}},
			},
			{
				GraphiteData.NullFloat64{sql.NullFloat64{Float64: 0.26, Valid: true}},
				GraphiteData.NullFloat64{sql.NullFloat64{Float64: 1.49660586e+09, Valid: true}},
			},
			{
				GraphiteData.NullFloat64{sql.NullFloat64{Float64: 0.3, Valid: true}},
				GraphiteData.NullFloat64{sql.NullFloat64{Float64: 1.49660592e+09, Valid: true}},
			},
			{
				GraphiteData.NullFloat64{sql.NullFloat64{Float64: 0.7, Valid: true}},
				GraphiteData.NullFloat64{sql.NullFloat64{Float64: 1.49660598e+09, Valid: true}},
			},
			{
				GraphiteData.NullFloat64{sql.NullFloat64{Float64: 0.5, Valid: true}},
				GraphiteData.NullFloat64{sql.NullFloat64{Float64: 1.49660604e+09, Valid: true}},
			},
			{
				GraphiteData.NullFloat64{sql.NullFloat64{Float64: 0.3, Valid: true}},
				GraphiteData.NullFloat64{sql.NullFloat64{Float64: 1.4966061e+09, Valid: true}},
			},
			{
				GraphiteData.NullFloat64{sql.NullFloat64{Float64: 0.25, Valid: true}},
				GraphiteData.NullFloat64{sql.NullFloat64{Float64: 1.49660616e+09, Valid: true}},
			},
			{
				GraphiteData.NullFloat64{sql.NullFloat64{Float64: 0, Valid: false}},
				GraphiteData.NullFloat64{sql.NullFloat64{Float64: 1.49660622e+09, Valid: true}},
			},
		},
	},
}

func TestCheckGraphiteNew(t *testing.T) {

	println("Testing 'CheckGraphite.New' function...")

	// Expected Result
	var expected = Alerts{
		Target:     "collectd.graphite.load.load.longterm",
		Datapoints: []float64{0.26, 0.3, 0.7, 0.5, 0.3, 0.25},
		Function:   "last",
		Scale:      2,
		Value:      0.25,
	}

	// Perform Function
	var value, _ = New(graphiteData[0], false, 2, "last")

	// Test Result
	var test = reflect.DeepEqual(value, expected)
	if test == false {
		t.Fatalf("Expected Value to be \"%v\" but was \"%v\"", expected, value)
	}

}

func TestCheckGraphiteNewZeros(t *testing.T) {

	println("Testing 'CheckGraphite.New' function with nulls as zeros...")

	// Expected Result
	var expected = Alerts{
		Target:     "collectd.graphite.load.load.longterm",
		Datapoints: []float64{0, 0, 0.26, 0.3, 0.7, 0.5, 0.3, 0.25},
		Function:   "last",
		Scale:      2,
		Value:      0.25,
	}

	// Perform Function
	var value, _ = New(graphiteData[0], true, 2, "last")

	// Test Result
	var test = reflect.DeepEqual(value, expected)
	if test == false {
		t.Fatalf("Expected Value to be \"%v\" but was \"%v\"", expected, value)
	}

}

func TestDoAlerts(t *testing.T) {

	println("Testing 'DoAlerts' method...")

	var value, _ = New(graphiteData[0], false, 2, "last")

	// OK: value less than Critical and Warning
	_, expected := value.DoAlerts(0.3, 0.4, false)
	println("Testing 'DoAlerts' OK status...")
	if expected != 0 {
		t.Fatalf("Expected Value to be \"%v\" but was \"%v\"", expected, value)
	}

	// Warning: value less than Critical, greater than Warning
	_, expected = value.DoAlerts(0.2, 0.3, false)
	println("Testing 'DoAlerts' Warning status...")
	if expected != 1 {
		t.Fatalf("Expected Value to be \"%v\" but was \"%v\"", expected, value)
	}

	// Critical: value greater than Critical, greater than Warning
	_, expected = value.DoAlerts(0.1, 0.2, false)
	println("Testing 'DoAlerts' Critical status...")
	if expected != 2 {
		t.Fatalf("Expected Value to be \"%v\" but was \"%v\"", expected, value)
	}

}

func TestDoAlertsInvert(t *testing.T) {

	println("Testing 'DoAlerts' method with inverted thresholds...")

	var value, _ = New(graphiteData[0], false, 2, "last")

	// OK: value greater than Critical and Warning
	_, expected := value.DoAlerts(0.2, 0.1, true)
	println("Testing 'DoAlerts' OK status...")
	if expected != 0 {
		t.Fatalf("Expected Value to be \"%v\" but was \"%v\"", expected, value)
	}

	// Warning: value greater than Critical, less than Warning
	_, expected = value.DoAlerts(0.3, 0.2, true)
	println("Testing 'DoAlerts' Warning status...")
	if expected != 1 {
		t.Fatalf("Expected Value to be \"%v\" but was \"%v\"", expected, value)
	}

	// Critical: value less than Critical, less than Warning
	_, expected = value.DoAlerts(0.4, 0.3, true)
	println("Testing 'DoAlerts' Critical status...")
	if expected != 2 {
		t.Fatalf("Expected Value to be \"%v\" but was \"%v\"", expected, value)
	}

}

func TestMinAggregation(t *testing.T) {

	println("Testing 'Min' aggregation method...")

	// Expected Result
	var expected = 0.25

	// Perform Function
	var value, _ = New(graphiteData[0], false, 2, "min")

	// Test Result
	if value.Value != expected {
		t.Fatalf("Expected Value to be \"%v\" but was \"%v\"", expected, value.Value)
	}

}

func TestMinAggregationZeros(t *testing.T) {

	println("Testing 'Min' aggregation method with nulls as zeros...")

	// Expected Result
	var expected = 0.0

	// Perform Function
	var value, _ = New(graphiteData[0], true, 2, "min")

	// Test Result
	if value.Value != expected {
		t.Fatalf("Expected Value to be \"%v\" but was \"%v\"", expected, value.Value)
	}

}

func TestMaxAggregation(t *testing.T) {

	println("Testing 'Max' aggregation method...")

	// Expected Result
	var expected = 0.7

	// Perform Function
	var value, _ = New(graphiteData[0], false, 2, "max")

	// Test Result
	if value.Value != expected {
		t.Fatalf("Expected Value to be \"%v\" but was \"%v\"", expected, value.Value)
	}

}

func TestMaxAggregationZeros(t *testing.T) {

	println("Testing 'Max' aggregation method with nulls as zeros...")

	// Expected Result
	var expected = 0.7

	// Perform Function
	var value, _ = New(graphiteData[0], true, 2, "max")

	// Test Result
	if value.Value != expected {
		t.Fatalf("Expected Value to be \"%v\" but was \"%v\"", expected, value.Value)
	}

}

func TestAvgAggregation(t *testing.T) {

	println("Testing 'Avg' aggregation method...")

	// Expected Result
	var expected = 0.39

	// Perform Function
	var value, _ = New(graphiteData[0], false, 2, "avg")

	// Test Result
	if value.Value != expected {
		t.Fatalf("Expected Value to be \"%v\" but was \"%v\"", expected, value.Value)
	}

}

func TestAvgAggregationZeros(t *testing.T) {

	println("Testing 'Avg' aggregation method with nulls as zeros...")

	// Expected Result
	var expected = 0.29

	// Perform Function
	var value, _ = New(graphiteData[0], true, 2, "avg")

	// Test Result
	if value.Value != expected {
		t.Fatalf("Expected Value to be \"%v\" but was \"%v\"", expected, value.Value)
	}

}

func TestSumAggregation(t *testing.T) {

	println("Testing 'Sum' aggregation method...")

	// Expected Result
	var expected = 2.31

	// Perform Function
	var value, _ = New(graphiteData[0], false, 2, "sum")

	// Test Result
	if value.Value != expected {
		t.Fatalf("Expected Value to be \"%v\" but was \"%v\"", expected, value.Value)
	}

}

func TestSumAggregationZeros(t *testing.T) {

	println("Testing 'Sum' aggregation method with nulls as zeros...")

	// Expected Result
	var expected = 2.31

	// Perform Function
	var value, _ = New(graphiteData[0], true, 2, "sum")

	// Test Result
	if value.Value != expected {
		t.Fatalf("Expected Value to be \"%v\" but was \"%v\"", expected, value.Value)
	}

}

func TestLastAggregation(t *testing.T) {

	println("Testing 'Last' aggregation method...")

	// Expected Result
	var expected = 0.25

	// Perform Function
	var value, _ = New(graphiteData[0], false, 2, "last")

	// Test Result
	if value.Value != expected {
		t.Fatalf("Expected Value to be \"%v\" but was \"%v\"", expected, value.Value)
	}

}

func TestLastAggregationZeros(t *testing.T) {

	println("Testing 'Last' aggregation method with nulls as zeros...")

	// Expected Result
	var expected = 0.25

	// Perform Function
	var value, _ = New(graphiteData[0], true, 2, "last")

	// Test Result
	if value.Value != expected {
		t.Fatalf("Expected Value to be \"%v\" but was \"%v\"", expected, value.Value)
	}

}
