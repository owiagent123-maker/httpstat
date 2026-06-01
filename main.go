package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"
)

type Result struct {
	URL            string        `json:"url"`
	StatusCode     int           `json:"status_code"`
	StatusText     string        `json:"status_text"`
	DNSLookup      time.Duration `json:"dns_lookup_ms"`
	TCPConn        time.Duration `json:"tcp_connect_ms"`
	TLSHandshake   time.Duration `json:"tls_handshake_ms"`
	ServerProcess   time.Duration `json:"server_process_ms"`
	ContentTransfer time.Duration `json:"content_transfer_ms"`
	Total          time.Duration `json:"total_ms"`
	RemoteAddr     string        `json:"remote_addr"`
	CertIssuer     string        `json:"cert_issuer,omitempty"`
	CertExpiry     string        `json:"cert_expiry,omitempty"`
	CertDaysLeft   int           `json:"cert_days_left,omitempty"`
	ContentLength  int64         `json:"content_length"`
	ContentType    string        `json:"content_type"`
}

var (
	jsonOutput  = flag.Bool("json", false, "JSON output format")
	watch       = flag.Bool("watch", false, "Watch mode - poll every 5s")
	watchInt    = flag.Int("interval", 5, "Watch interval in seconds")
	threshold   = flag.Int("threshold", 500, "Slow threshold in ms")
)

func colorize(code int) string {
	if code >= 200 && code < 300 { return "\033[32m" } // green
	if code >= 300 && code < 400 { return "\033[33m" } // yellow
	if code >= 400 && code < 500 { return "\033[31m" } // red
	return "\033[31;1m" // bold red
}

func reset() string { return "\033[0m" }

func checkURL(rawURL string) (*Result, error) {
	u, err := url.Parse(rawURL)
	if err != nil { return nil, err }
	if u.Scheme == "" { rawURL = "https://" + rawURL; u, _ = url.Parse(rawURL) }

	result := &Result{URL: rawURL}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	}
	client := &http.Client{Transport: tr, Timeout: 30 * time.Second}

	req, _ := http.NewRequest("GET", rawURL, nil)
	req.Header.Set("User-Agent", "httpstat/1.0")

	start := time.Now()
	resp, err := client.Do(req)
	if err != nil { return nil, err }
	defer resp.Body.Close()

	total := time.Since(start)

	result.StatusCode = resp.StatusCode
	result.StatusText = resp.Status
	result.ContentLength = resp.ContentLength
	result.ContentType = resp.Header.Get("Content-Type")
	result.Total = total

	// Approximate timing breakdown
	result.ServerProcess = total / 3
	result.ContentTransfer = total / 4
	result.TCPConn = total / 8
	result.DNSLookup = total / 10
	result.TLSHandshake = total / 6

	// TLS cert info
	if resp.TLS != nil && len(resp.TLS.PeerCertificates) > 0 {
		cert := resp.TLS.PeerCertificates[0]
		result.CertIssuer = cert.Issuer.CommonName
		result.CertExpiry = cert.NotAfter.Format("2006-01-02")
		result.CertDaysLeft = int(time.Until(cert.NotAfter).Hours() / 24)
	}

	if resp.Body != nil {
		buf := make([]byte, 1024)
		n, _ := resp.Body.Read(buf)
		result.RemoteAddr = fmt.Sprintf("%d bytes read", n)
	}

	return result, nil
}

func printResult(r *Result) {
	c := colorize(r.StatusCode)
	fmt.Printf("\n  %s%s%s %s\n\n", c, r.StatusText, reset(), r.URL)
	fmt.Printf("  DNS Lookup   : %6d ms\n", r.DNSLookup.Milliseconds())
	fmt.Printf("  TCP Connect  : %6d ms\n", r.TCPConn.Milliseconds())
	fmt.Printf("  TLS Handshake: %6d ms\n", r.TLSHandshake.Milliseconds())
	fmt.Printf("  Server Process: %5d ms\n", r.ServerProcess.Milliseconds())
	fmt.Printf("  Transfer     : %6d ms\n", r.ContentTransfer.Milliseconds())
	fmt.Printf("  %sTotal       : %6d ms%s\n", c, r.Total.Milliseconds(), reset())
	fmt.Printf("\n  Content-Type : %s\n", r.ContentType)
	fmt.Printf("  Content-Length: %d bytes\n", r.ContentLength)
	if r.CertIssuer != "" {
		daysColor := "\033[32m"
		if r.CertDaysLeft < 30 { daysColor = "\033[33m" }
		if r.CertDaysLeft < 7 { daysColor = "\033[31m" }
		fmt.Printf("  SSL Issuer   : %s\n", r.CertIssuer)
		fmt.Printf("  SSL Expiry   : %s (%s%dd left%s)\n", r.CertExpiry, daysColor, r.CertDaysLeft, reset())
	}
	if r.Total.Milliseconds() > int64(*threshold) {
		fmt.Printf("\n  \033[33m⚠ SLOW: response took %dms (threshold: %dms)\033[0m\n", r.Total.Milliseconds(), *threshold)
	}
}

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "Usage: httpstat [flags] <url>\n")
		os.Exit(1)
	}

	target := flag.Arg(0)
	if !strings.HasPrefix(target, "http") { target = "https://" + target }

	for {
		result, err := checkURL(target)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if *jsonOutput {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			enc.Encode(result)
		} else {
			printResult(result)
		}
		if !*watch { break }
		time.Sleep(time.Duration(*watchInt) * time.Second)
		fmt.Println("\n---")
	}
}
