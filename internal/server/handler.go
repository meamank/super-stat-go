package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/meamank/super-stat-go/internal/cpu"
	"github.com/meamank/super-stat-go/internal/mem"
)

type StreamPayload struct {
	CPU map[string]float64 `json:"cpu"`
	Mem mem.MemInfo        `json:"mem"`
}

func CPUStatStreamHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)

	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	for {
		select {
		case <-r.Context().Done():
			fmt.Println("Client Disconnected from stream!")
			return
		default:
			cpuUsage, err := cpu.GetPerCPUCore()

			if err != nil {
				continue
			}

			memStat, err := mem.GetMemUsage()

			if err != nil {
				continue
			}

			payload := StreamPayload{
				CPU: cpuUsage,
				Mem: memStat,
			}

			dataBytes, err := json.Marshal(payload)

			if err != nil {
				continue
			}

			fmt.Fprintf(w, "data: %s\n\n", dataBytes)
			flusher.Flush()
		}
	}

}
