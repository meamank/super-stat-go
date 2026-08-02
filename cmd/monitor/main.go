package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/meamank/super-stat-go/internal/cpu"
)

func main() {

	usage, err := cpu.GetPerCPUCore()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	var keys []string

	for k := range usage {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		fmt.Printf("%-8s: %.2f%%\n", strings.ToUpper(k), usage[k])
	}

}
