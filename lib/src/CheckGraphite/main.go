package CheckGraphite

import (
	"GraphiteData"
	"fmt"
	"math"
)

// Alerts struct
type Alerts struct {
	Target     string
	Datapoints []float64
	Function   string
	Value      float64
}

// New function aggregates Graphite data for alerting
func New(metrics GraphiteData.GraphiteMetrics, zero bool, fctn string) (m Alerts, err error) {

	// Set Target
	m.Target = metrics.Target

	// Add Datapoints to CheckGraphite data
	for _, d := range metrics.Datapoints {
		if zero {
			// Skip the last Datapoint if it is a null
			if len(m.Datapoints) == len(metrics.Datapoints)-1 && d[0].Valid == false {
				break
			}
			// Convert null to 0
			m.Datapoints = append(m.Datapoints, RoundAsFloat64(d[0].Float64, 2))
		} else {
			// Add only float64 values
			if d[0].Valid {
				m.Datapoints = append(m.Datapoints, RoundAsFloat64(d[0].Float64, 2))
			}
		}
	}

	// Set Aggregation Function
	m.Function = fctn

	// Set Aggregated Value. Default is 'last'.
	if fctn == "min" {
		m.Value = m.Min()
	} else if fctn == "max" {
		m.Value = m.Max()
	} else if fctn == "avg" {
		m.Value = m.Avg()
	} else if fctn == "sum" {
		m.Value = m.Sum()
	} else {
		m.Value = m.Last()
	}

	// Return CheckGraphite data
	return m, nil
}

// DoAlerts method
func (m *Alerts) DoAlerts(warning float64, critical float64, invert bool) (msg string, exit int) {

	// Interpret Thresholds
	if invert == false {
		if m.Value >= critical {
			msg = "CRITICAL: "
			exit = 2
		} else if m.Value >= warning {
			msg = "WARNING: "
			exit = 1
		} else {
			msg = "OK: "
			exit = 0
		}
	} else {
		if m.Value <= critical {
			msg = "CRITICAL: "
			exit = 2
		} else if m.Value <= warning {
			msg = "WARNING: "
			exit = 1
		} else {
			msg = "OK: "
			exit = 0
		}
	}

	// Message data
	msg += fmt.Sprintf("Target: %s, ValueCount: %d, Function: %s, Value: %.2f",
		m.Target,
		len(m.Datapoints),
		m.Function,
		m.Value,
	)

	// Return Message and Exit Code
	return msg, exit
}

// Min Aggregation method
func (m *Alerts) Min() (min float64) {
	min = m.Datapoints[0]
	for i := 1; i < len(m.Datapoints); i++ {
		if m.Datapoints[i] < min {
			min = m.Datapoints[i]
		}
	}
	return RoundAsFloat64(min, 2)
}

// Max Aggregation method
func (m *Alerts) Max() (max float64) {
	for i := 0; i < len(m.Datapoints); i++ {
		if m.Datapoints[i] > max {
			max = m.Datapoints[i]
		}
	}
	return RoundAsFloat64(max, 2)
}

// Avg Aggregation method
func (m *Alerts) Avg() (avg float64) {
	for _, v := range m.Datapoints {
		avg += v
	}
	return RoundAsFloat64(avg/float64(len(m.Datapoints)), 2)
}

// Sum Aggregation method
func (m *Alerts) Sum() (sum float64) {
	for _, v := range m.Datapoints {
		sum += v
	}
	return RoundAsFloat64(sum, 2)
}

// Last Aggregation method
func (m *Alerts) Last() float64 {
	return RoundAsFloat64(m.Datapoints[len(m.Datapoints)-1], 2)
}

// RoundAsInt function
func RoundAsInt(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

// RoundAsFloat64 function: Specify Float64 Precision
func RoundAsFloat64(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(RoundAsInt(num*output)) / output
}
