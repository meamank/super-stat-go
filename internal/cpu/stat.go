package cpu

import (
	"bufio"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

type CPUStat struct {
	total uint64
	idle  uint64
}

func parseCPUStatFromReader(r io.Reader) (map[string]CPUStat, error) {
	statsMap := make(map[string]CPUStat)
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()

		if !strings.HasPrefix(line, "cpu") {

			break
		}

		fields := strings.Fields(line)

		coreName := fields[0]

		if coreName == "cpu" {
			coreName = "total"
		}

		var total, idle uint64

		for i, valStr := range fields[1:] {
			valInt, _ := strconv.ParseUint(valStr, 10, 64)

			total += valInt

			if i == 3 || i == 4 {
				idle += valInt
			}

		}

		statsMap[coreName] = CPUStat{total: total, idle: idle}

	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return statsMap, nil
}

func readAllCPUStat() (map[string]CPUStat, error) {
	file, err := os.Open("/proc/stat")

	if err != nil {
		return nil, err
	}

	defer file.Close()

	return parseCPUStatFromReader(file)
}

func GetPerCPUCore() (map[string]float64, error) {
	firstStat, err := readAllCPUStat()

	if err != nil {
		return nil, err
	}

	time.Sleep(1 * time.Second)

	secondStat, err := readAllCPUStat()

	if err != nil {
		return nil, err
	}

	results := make(map[string]float64)

	for coreName, secondStat := range secondStat {
		firstStat, exists := firstStat[coreName]

		if !exists {
			continue
		}

		totalDelta := float64(secondStat.total - firstStat.total)
		idleDelta := float64(secondStat.idle - firstStat.idle)

		if totalDelta > 0 {
			usage := ((totalDelta - idleDelta) / totalDelta) * 100
			results[coreName] = math.Round(usage*100) / 100
		} else {
			results[coreName] = 0.0
		}
	}
	return results, nil
}
