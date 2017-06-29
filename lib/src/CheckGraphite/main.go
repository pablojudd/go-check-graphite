package CheckGraphite

import (
	"fmt"
	"math"

	"github.com/pablojudd/go-graphite-getmetrics"
)

// Alerts struct
type Alerts struct {
	Target     string
	Datapoints []float64
	Function   string
	Scale      int
	Value      float64
}

// New function aggregates Graphite data for alerting
func New(metrics GraphiteData.GraphiteMetrics, zero bool, scale int, fctn string) (m Alerts, err error) {

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
			m.Datapoints = append(m.Datapoints, roundAsFloat64(d[0].Float64, scale))
		} else {
			// Add only float64 values
			if d[0].Valid {
				m.Datapoints = append(m.Datapoints, roundAsFloat64(d[0].Float64, scale))
			}
		}
	}

	// Set Aggregation Function
	m.Function = fctn

	// Set Numeric Scale
	m.Scale = scale

	// Set Aggregated Value. Default is 'last'.
	if fctn == "min" {
		m.Value = m.min()
	} else if fctn == "max" {
		m.Value = m.max()
	} else if fctn == "avg" {
		m.Value = m.avg()
	} else if fctn == "sum" {
		m.Value = m.sum()
	} else {
		m.Value = m.last()
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

// Min Aggregation
func (m *Alerts) min() float64 {
	min := m.Datapoints[0]
	for _, v := range m.Datapoints[1:] {
		if v < min {
			min = v
		}
	}
	return roundAsFloat64(min, m.Scale)
}

// Max Aggregation
func (m *Alerts) max() float64 {
	max := m.Datapoints[0]
	for _, v := range m.Datapoints[1:] {
		if v > max {
			max = v
		}
	}
	return roundAsFloat64(max, m.Scale)
}

// Avg Aggregation
func (m *Alerts) avg() float64 {
	avg := 0.0
	for _, v := range m.Datapoints {
		avg += v
	}
	avg = avg / float64(len(m.Datapoints))
	return roundAsFloat64(avg, m.Scale)
}

// Sum Aggregation
func (m *Alerts) sum() float64 {
	sum := 0.0
	for _, v := range m.Datapoints {
		sum += v
	}
	return roundAsFloat64(sum, m.Scale)
}

// Last Aggregation
func (m *Alerts) last() float64 {
	last := m.Datapoints[len(m.Datapoints)-1]
	return roundAsFloat64(last, m.Scale)
}

// roundAsInt function
func roundAsInt(num float64) int {
	// Use math.Copysign to handle negative num values
	return int(num + math.Copysign(0.5, num))
}

// roundAsFloat64 function: Specify the numeric scale for a float64 value
func roundAsFloat64(num float64, scale int) float64 {
	output := math.Pow(10, float64(scale))
	return float64(roundAsInt(num*output)) / output
}
