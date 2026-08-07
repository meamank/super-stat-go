# `ssgo` — Super Stat Go

A lightweight, zero-dependency Linux CPU monitoring CLI tool and endpoints for web built with Go.

---

## Features

- **Zero External Dependencies:** Uses Go's standard packages.
- **Per-Core & Aggregated Breakdown:** Calculates overall CPU usage percentage as well as individual core metrics.
- **Memory & Swap Telemetry:** Track physical RAM usage and swap capacity.
- **Live Terminal Watch Mode (`--watch`):** Terminal watch mode for continuous stat.
- **Web Server Mode (`--server`):** Pushes real-time JSON metrics over HTTP (`GET /cpu-stat`) to feed web dashboards.
- **Single-Binary Dual Mode:** Operates seamlessly as both a CLI tool and an HTTP server from a single binary.

---

## Installation

### Option 1: Via `go install` (Recommended)

Requires Go 1.21+ installed on a Linux system:

```bash
go install github.com/meamank/super-stat-go/cmd/ssgo@latest
```

Ensure `$(go env GOPATH)/bin` is in your `$PATH`.

---

### Option 2: Building from Source

```bash
git clone https://github.com/meamank/super-stat-go.git
cd super-stat-go

# Build binary
go build -o ssgo ./cmd/ssgo

# (Optional) Move to system PATH
sudo mv ssgo /usr/local/bin/
```

---

### Option 3: Cross-Compiling for Linux from macOS

If you are developing on macOS for a home Linux server / Raspberry Pi:

```bash
# For x86_64 Linux Server:
GOOS=linux GOARCH=amd64 go build -o ssgo ./cmd/ssgo

# For ARM64 Linux / Raspberry Pi:
GOOS=linux GOARCH=arm64 go build -o ssgo ./cmd/ssgo
```

Deploy the binary to your server:
```bash
scp ssgo user@your-server-ip:/tmp/
ssh user@your-server-ip "sudo mv /tmp/ssgo /usr/local/bin/ && sudo chmod +x /usr/local/bin/sudo"
```

---

## Usage

### 1. One-Time CPU Stat
Print current CPU usage once and exit:

```bash
ssgo
```

**Output:**
```text
TOTAL   : 1.25%
CPU0    : 1.80%
CPU1    : 0.70%
```

---

### 2. Live Terminal Watch Mode
Stream continuous per-core CPU stats with in-place terminal redrawing (Ctrl+C to stop):

```bash
ssgo --watch
```

---

### 3. Web Streaming Server Mode (for React / Web Dashboards)
Start an HTTP server on port `:8080`:

```bash
ssgo --server
```

Open `http://localhost:8080/cpu-stat` in your browser to view the live JSON stream:

```text
data: {"cpu":{"cpu0":0,"cpu1":1.00,"cpu2":0.99,"total":0.76},"mem":{"totalMB":2767.67,"availableMB":2133.36,"usedMB":634.31,"usagePercent":22.92,"swapTotalMB":2922.00,"swapUsedMB":0.00,"swapUsagePercent":0.00}}
```

#### React Integration Example:
```javascript
useEffect(() => {
  const eventSource = new EventSource('http://your-server:8080/cpu-stat');
  eventSource.onmessage = (event) => {
    const data = JSON.parse(event.data);
    console.log("Telemetry Payload:", data.cpu, data.mem);
  };
  return () => eventSource.close();
}, []);
```

---

### Memory Telemetry

Tracks physical RAM and Swap Memory using the standard Linux kernel formula:
- **Used RAM:** `MemTotal - MemAvailable`
- **RAM Usage %:** `(MemTotal - MemAvailable) / MemTotal * 100`
- **Swap Used:** `SwapTotal - SwapFree`

---

### Project Architecture

Organized according to standard Go project layout:

```text
super-stat-go/
├── cmd/
│   └── ssgo/
│       └── main.go         # Entrypoint & CLI flag routing
├── internal/
│   ├── cpu/
│   │   └── stat.go         # Linux /proc/stat parser
│   ├── mem/
│   │   └── stat.go         # Linux /proc/meminfo parser
│   └── server/
│       └── handler.go      # HTTP Server-Sent Events (SSE) handler
├── embed.go                # Static asset embedding directive
└── go.mod
```

---

### Running via Docker (for local testing on macOS)

```bash
docker run --rm -p 8080:8080 -v $(pwd):/app -w /app golang:1.26 go run ./cmd/ssgo --server
```

---

## License

MIT © [Aman](https://github.com/meamank)
