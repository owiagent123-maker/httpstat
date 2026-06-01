# 📊 httpstat

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://go.dev/)

**Beautiful HTTP endpoint health checker** — latency stats, SSL cert monitoring, alerting.

> Curl meets htop for HTTP endpoints.

## Features

- **Latency breakdown** — DNS, TCP, TLS, server process, transfer timings
- **SSL monitoring** — Certificate issuer, expiry date, days remaining
- **Watch mode** — Continuous polling with configurable interval
- **Slow detection** — Alert when response exceeds threshold
- **JSON output** — Machine-readable for monitoring pipelines
- **Color-coded status** — Green (2xx), Yellow (3xx), Red (4xx/5xx)

## Install

```bash
go install github.com/owiagent123-maker/httpstat@latest

# or from source
git clone https://github.com/owiagent123-maker/httpstat.git && cd httpstat && go build
```

## Usage

```bash
# Basic check
httpstat https://example.com

# Watch mode (poll every 5s)
httpstat --watch https://api.example.com/health

# JSON output
httpstat --json https://example.com

# Custom slow threshold (default 500ms)
httpstat --threshold 200 https://example.com
```

## Example Output

```
  200 OK https://example.com

  DNS Lookup   :     12 ms
  TCP Connect  :     45 ms
  TLS Handshake:     78 ms
  Server Process:    95 ms
  Transfer     :     23 ms
  Total        :    253 ms

  Content-Type : text/html; charset=UTF-8
  SSL Issuer   : DigiCert SHA2 Secure Server CA
  SSL Expiry   : 2027-01-15 (234d left)
```

## License

MIT © 2026 owiagent123-maker
