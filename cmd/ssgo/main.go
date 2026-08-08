package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"
	"time"

	superstatgo "github.com/meamank/super-stat-go"
	"github.com/meamank/super-stat-go/internal/cpu"
	"github.com/meamank/super-stat-go/internal/mem"
	"github.com/meamank/super-stat-go/internal/server"
)

func formatUsage(payload server.StreamPayload) {
	fmt.Print("\033[H\033[2J") // In-place clear terminal

	fmt.Println("==========================================")
	fmt.Println("    🚀 SUPER STAT MONITOR (Ctrl+C to exit)")
	fmt.Println("==========================================")

	// 1. CPU Telemetry
	fmt.Println("--- CPU TELEMETRY ---")
	var keys []string
	for k := range payload.CPU {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Printf("%-8s: %.2f%%\n", strings.ToUpper(k), payload.CPU[k])
	}

	// 2. Memory Telemetry
	fmt.Println("\n--- MEMORY TELEMETRY ---")
	fmt.Printf("%-12s: %.2f MB / %.2f MB (%.2f%%)\n", "RAM USED", payload.Mem.UsedMB, payload.Mem.
		TotalMB, payload.Mem.UsagePercent)
	fmt.Printf("%-12s: %.2f MB / %.2f MB (%.2f%%)\n", "SWAP USED", payload.Mem.SwapUsedMB, payload.
		Mem.SwapTotalMB, payload.Mem.SwapUsagePercent)
	fmt.Println("==========================================")
}

func main() {

	watchMode := flag.Bool("watch", false, "Stream CPU usage!")
	serverMode := flag.Bool("server", false, "HTTP Endpoint /cpu-stat")

	flag.Parse()

	if *serverMode {
		mux := http.NewServeMux()

		mux.HandleFunc("GET /cpu-stat", server.CPUStatStreamHandler)

		fileServer := http.FileServer(superstatgo.GetFileSystem())
		mux.Handle("/", fileServer)
		fmt.Println("server running on Port :8080")
		log.Fatal(http.ListenAndServe(":8080", mux))
		return
	}

	if *watchMode {
		ticker := time.NewTicker(1 * time.Second)

		defer ticker.Stop()

		sigChan := make(chan os.Signal, 1)

		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

		for {
			select {
			case <-ticker.C:
				cpuUsage, err := cpu.GetPerCPUCore()

				if err != nil {
					fmt.Fprintf(os.Stderr, "Error: %v\n", err)
					os.Exit(1)
				}
				memUsage, err := mem.GetMemUsage()

				if err != nil {
					fmt.Fprintf(os.Stderr, "Error: %v\n", err)
					os.Exit(1)
				}

				payload := server.StreamPayload{
					CPU: cpuUsage,
					Mem: memUsage,
				}

				formatUsage(payload)
			case <-sigChan:
				return
			}

		}
	}
	cpuUsage, err := cpu.GetPerCPUCore()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	memUsage, err := mem.GetMemUsage()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	payload := server.StreamPayload{
		CPU: cpuUsage,
		Mem: memUsage,
	}

	formatUsage(payload)

}
