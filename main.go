package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

type Result struct {
	URL           string `json:"url"`
	Status        int    `json:"status"`
	StatusText    string `json:"status_text"`
	TotalMs       int64  `json:"total_ms"`
	ContentLength int64  `json:"content_length"`
	SSLExpiry     string `json:"ssl_expiry,omitempty"`
	SSLIssuer     string `json:"ssl_issuer,omitempty"`
	HasSSL        bool   `json:"has_ssl"`
}

func check(url string, showSSL bool) Result {
	start := time.Now()
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(url)
	total := time.Since(start)

	r := Result{URL: url, TotalMs: total.Milliseconds()}
	if err != nil {
		r.StatusText = fmt.Sprintf("Error: %v", err)
		return r
	}
	defer resp.Body.Close()

	r.Status = resp.StatusCode
	r.StatusText = resp.Status
	r.ContentLength = resp.ContentLength

	if showSSL && resp.TLS != nil && len(resp.TLS.PeerCertificates) > 0 {
		cert := resp.TLS.PeerCertificates[0]
		r.SSLExpiry = cert.NotAfter.Format("2006-01-02 15:04:05")
		r.SSLIssuer = cert.Issuer.CommonName
		r.HasSSL = true
	}

	return r
}

func hostFromURL(u string) string {
	u = strings.TrimPrefix(u, "https://")
	u = strings.TrimPrefix(u, "http://")
	if i := strings.Index(u, "/"); i > 0 {
		u = u[:i]
	}
	return u
}

func printResult(r Result) {
	color := "\033[32m"
	if r.Status >= 400 || r.Status == 0 {
		color = "\033[31m"
	} else if r.Status >= 300 {
		color = "\033[33m"
	}
	reset := "\033[0m"

	fmt.Printf("\n\033[1mhttpstat:\033[0m %s%s %d %s%s\n", color, r.URL, r.Status, r.StatusText, reset)
	fmt.Printf("  Total:     \033[1m%d ms\033[0m\n", r.TotalMs)
	if r.ContentLength > 0 {
		fmt.Printf("  Size:      %d bytes\n", r.ContentLength)
	}
	if r.HasSSL {
		fmt.Printf("  SSL:       %s (Issuer: %s)\n", r.SSLExpiry, r.SSLIssuer)
	}
	fmt.Println()
}

func main() {
	jsonOut := flag.Bool("json", false, "JSON output")
	showSSL := flag.Bool("ssl", true, "Show SSL cert info")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "httpstat — Beautiful HTTP health checker\n\nUsage: httpstat [flags] <URL> [URL...]\n\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	var results []Result
	for _, u := range flag.Args() {
		if !strings.HasPrefix(u, "http") {
			u = "https://" + u
		}
		r := check(u, *showSSL)
		results = append(results, r)
		if !*jsonOut {
			printResult(r)
		}
	}

	if *jsonOut {
		data, _ := json.MarshalIndent(results, "", "  ")
		fmt.Println(string(data))
	}

	_ = net.Dialer{}
}
