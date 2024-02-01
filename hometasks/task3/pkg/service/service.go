package service

import (
	"fmt"
	"runtime"

	"github.com/shirou/gopsutil/cpu"
)

type Metrics struct{}

func NewMetrics() *Metrics {
	return &Metrics{}
}

func (m *Metrics) GetCPULoad() (float64, error) {
	percentages, err := cpu.Percent(0, false)
	if err != nil {
		return 0, fmt.Errorf("GetCPULoad: cpu.Percent() %v", err)
	}

	if len(percentages) > 0 {
		return percentages[0], nil
	}
	return 0, fmt.Errorf("сould not get CPU load information")
}
func (m *Metrics) GetThreadCount() int {
	return runtime.NumCPU()
}
