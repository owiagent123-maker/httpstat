# ⚡ httpstat

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org/)

**Beautiful HTTP endpoint health checker** with latency stats and SSL monitoring.

> Curl meets htop for HTTP endpoints.

## Install

```bash
go install github.com/owiagent123-maker/httpstat@latest

# Or from source
git clone https://github.com/owiagent123-maker/httpstat.git
cd httpstat
go build -o httpstat .
```

## Usage

```bash
# Single endpoint
httpstat https://example.com

# Multiple endpoints
httpstat https://google.com https://github.com https://cloudflare.com

# JSON output
httpstat --json https://api.github.com

# Without SSL check
httpstat --ssl=false http://localhost:3000
```

## Sample Output

```
httpstat: https://example.com 200 OK
  Total:     142 ms
  Size:      1256 bytes
  SSL:       2027-01-15 23:59:59 (Issuer: DigiCert SHA2 Secure)
```

## Features

- **Latency tracking** — Total response time in milliseconds
- **SSL monitoring** — Certificate expiry and issuer info
- **Multi-endpoint** — Check multiple URLs in one command
- **JSON output** — Machine-readable for monitoring pipelines
- **Color-coded** — Green (2xx), Yellow (3xx), Red (4xx/5xx)

## License

MIT © 2026 owiagent123-maker
