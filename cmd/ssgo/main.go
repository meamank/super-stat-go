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

	"github.com/meamank/super-stat-go/internal/cpu"
	"github.com/meamank/super-stat-go/internal/server"
)

func formatUsage(usage map[string]float64) {
	fmt.Print("\033[H\033[2J")
	var keys []string

	for k := range usage {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		fmt.Printf("%-8s: %.2f%%\n", strings.ToUpper(k), usage[k])
	}
}

func main() {

	watchMode := flag.Bool("watch", false, "Stream CPU usage!")
	serverMode := flag.Bool("server", false, "HTTP Endpoint /cpu-stat")

	flag.Parse()

	if *serverMode {
		mux := http.NewServeMux()

		mux.HandleFunc("GET /cpu-stat", server.CPUStatStreamHandler)

		fileServer := http.FileServer(server.GetFileSystem())
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
				usage, err := cpu.GetPerCPUCore()

				if err != nil {
					fmt.Fprintf(os.Stderr, "Error: %v\n", err)
					os.Exit(1)
				}
				formatUsage(usage)
			case <-sigChan:
				return
			}

		}
	}
	usage, err := cpu.GetPerCPUCore()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	formatUsage(usage)

}
