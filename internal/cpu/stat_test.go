package cpu

import (
	"strings"
	"testing"
)

func TestParseCPUStatFromReader(t *testing.T) {
	tests := []struct {
		name        string
		sampleInput string
		wantCores   []string
		wantTotal   uint64
		wantIdle    uint64
		wantErr     bool
	}{
		{
			name: "valid multi-core sample",
			sampleInput: `cpu  1000 100 500 8000 200 0 0 0 0 0
cpu0 500 50 250 4000 100 0 0 0 0 0
cpu1 500 50 250 4000 100 0 0 0 0 0
intr 123456`,
			wantCores: []string{"total", "cpu0", "cpu1"},
			wantTotal: 9800, // 1000+100+500+8000+200
			wantIdle:  8200, // 8000 (idle) + 200 (iowait)
			wantErr:   false,
		},
		{
			name:        "non-cpu content",
			sampleInput: "memory 123456\nintr 789",
			wantCores:   []string{},
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.sampleInput)

			stats, err := parseCPUStatFromReader(reader)

			if (err != nil) != tt.wantErr {
				t.Fatalf("wantErr %v, got err: %v", tt.wantErr, err)
			}

			for _, core := range tt.wantCores {
				if _, exists := stats[core]; !exists {
					t.Errorf("Expected core %s to exist in stats map", core)
				}
			}

			if len(tt.wantCores) > 0 {
				if totalStat, exists := stats["total"]; exists {
					if totalStat.total != tt.wantTotal {
						t.Errorf("total ticks got %d, want %d", totalStat.total, tt.wantTotal)
					}
					if totalStat.idle != tt.wantIdle {
						t.Errorf("idle ticks got %d, want %d", totalStat.idle, tt.wantIdle)
					}
				}
			}
		})
	}
}
