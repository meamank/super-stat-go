package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/meamank/super-stat-go/internal/cpu"
)

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
			usage, err := cpu.GetPerCPUCore()

			if err != nil {
				continue
			}

			dataBytes, err := json.Marshal(usage)

			if err != nil {
				continue
			}

			fmt.Fprintf(w, "data: %s\n\n", dataBytes)
			flusher.Flush()
		}
	}

}
