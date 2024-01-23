package service

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"runtime"
)

func GetCPULoad() (float64, error) {
	percentages, err := cpu.Percent(0, false)
	if err != nil {
		return 0, err
	}

	if len(percentages) > 0 {
		return percentages[0], nil
	}

	return 0, fmt.Errorf("сould not get CPU load information")
}
func GetThreadCount() int {
	return runtime.NumCPU()
}
