# `ssgo` — Super Stat Go

A lightweight, zero-dependency Linux CPU monitoring CLI tool and endpoints for web built with Go.

---

## Features

- **Zero External Dependencies:** Uses Go's standard packages.
- **Per-Core & Aggregated Breakdown:** Calculates overall CPU usage percentage as well as individual core metrics.
- **Live Terminal Watch Mode (`--watch`):** Terminal watch mode for continuous stat.
- **Web Server Mode (`--server`):** Pushes real-time JSON metrics over HTTP (`GET /cpu-stream`) to feed web dashboards.
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
ssh user@your-server-ip "sudo mv /tmp/ssgo /usr/local/bin/ && sudo chmod +x /usr/local/bin/ssgo"
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

Open `http://localhost:8080/cpu-stream` in your browser to view the live JSON stream:

```text
data: {"timestamp":"16:31:18","usage":1.25}

data: {"timestamp":"16:31:19","usage":0.80}
```

#### React Integration Example:
```javascript
useEffect(() => {
  const eventSource = new EventSource('http://your-server:8080/cpu-stream');
  eventSource.onmessage = (event) => {
    const data = JSON.parse(event.data);
    console.log("New CPU point:", data); // { timestamp: "16:31:18", usage: 1.25 }
  };
  return () => eventSource.close();
}, []);
```

---

### Running via Docker (for local testing on macOS)

```bash
docker run --rm -p 8080:8080 -v $(pwd):/app -w /app golang:1.26 go run ./cmd/ssgo --server
```

---



## License

MIT © [Aman](https://github.com/meamank)
