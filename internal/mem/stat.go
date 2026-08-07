package mem

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type MemInfo struct {
	TotalMB      float64 `json:"totalMB"`
	AvailableMB  float64 `json:"availableMB"`
	UsedMB       float64 `json:"usedMB"`
	UsagePercent float64 `json:"usagePercent"`

	SwapTotalMB      float64 `json:"swapTotalMB"`
	SwapUsedMB       float64 `json:"swapUsedMB"`
	SwapUsagePercent float64 `json:"swapUsagePercent"`
}

func GetMemUsage() (MemInfo, error) {
	file, err := os.Open("/proc/meminfo")

	if err != nil {
		return MemInfo{}, err
	}

	defer file.Close()

	var memTotal, memAvailable, swapTotal, swapFree uint64
	var foundTotal, foundAvail bool

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)

		if len(fields) < 2 {
			continue
		}

		key := fields[0]

		val, _ := strconv.ParseUint(fields[1], 10, 64)

		switch key {
		case "MemTotal:":
			memTotal = val
			foundTotal = true
		case "MemAvailable:":
			memAvailable = val
			foundAvail = true
		case "SwapTotal:":
			swapTotal = val
		case "SwapFree:":
			swapFree = val
		}

	}

	if err := scanner.Err(); err != nil {
		return MemInfo{}, err
	}

	if !foundTotal || !foundAvail {
		return MemInfo{}, fmt.Errorf("Missing Required meminfo fields")
	}

	totalRamMB := float64(memTotal) / 1024
	availRamMB := float64(memAvailable) / 1024
	usedRamMB := totalRamMB - availRamMB
	usageRamPercent := (usedRamMB / totalRamMB) * 100

	swapTotalMB := float64(swapTotal) / 1024
	swapFreeMB := float64(swapFree) / 1024
	swapUsedMB := swapTotalMB - swapFreeMB

	var swapPercent float64

	if swapTotalMB > 0 {
		swapPercent = (swapUsedMB / swapTotalMB) * 100
	}

	return MemInfo{
		TotalMB:          math.Round(totalRamMB*100) / 100,
		UsedMB:           math.Round(usedRamMB*100) / 100,
		AvailableMB:      math.Round(availRamMB*100) / 100,
		UsagePercent:     math.Round(usageRamPercent*100) / 100,
		SwapTotalMB:      math.Round(swapTotalMB*100) / 100,
		SwapUsedMB:       math.Round(swapUsedMB*100) / 100,
		SwapUsagePercent: math.Round(swapPercent*100) / 100,
	}, nil

}
