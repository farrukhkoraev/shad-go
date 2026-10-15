//go:build !solution

package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"flag"
	
	"net/url"
	"net/http"
	"net/http/httputil"


	//	"golang.org/x/text/cases"
	"gopkg.in/yaml.v3"
)

type Rule struct {
    Endpoint                string   `yaml:"endpoint"`
    ForbiddenUserAgents     []string `yaml:"forbidden_user_agents"`
    ForbiddenHeaders        []string `yaml:"forbidden_headers"`
    RequiredHeaders         []string `yaml:"required_headers"`
    MaxRequestLengthBytes   int      `yaml:"max_request_length_bytes"`
    MaxResponseLengthBytes  int      `yaml:"max_response_length_bytes"`
    ForbiddenResponseCodes  []int    `yaml:"forbidden_response_codes"`
    ForbiddenRequestRe      []string `yaml:"forbidden_request_re"`
    ForbiddenResponseRe     []string `yaml:"forbidden_response_re"`
}

type Config struct {
	Rules []Rule
}

type Firewall struct {
	transport http.RoundTripper
	rules map[string]Rule
}

func NewFirewall (config Config) *Firewall {
	rules := map[string]Rule {}
	for _, rule := range config.Rules {
		rules[rule.Endpoint] = rule
	}

	return &Firewall { transport: http.DefaultTransport, rules: rules }
}


func (fw *Firewall) RoundTrip(r *http.Request) (*http.Response, error) {
	
	rule, ok := fw.rules[r.URL.Path]

	if !ok {
		return fw.transport.RoundTrip(r)
	}
	body, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewReader(body))
	
	if !allowRequest(r.Header, body, rule) {
		return forbidden(), nil
	}

	resp, err := fw.transport.RoundTrip(r)
	if err != nil {
		return resp, err
	}
	body, _ = io.ReadAll(resp.Body)
	resp.Body = io.NopCloser(bytes.NewReader(body))
	if !allowResponse(resp.StatusCode, body, rule) {
		return forbidden(), nil
	}

	return resp, err
}

func main() {
	addr := flag.String("addr", "localhost:8080", "listen address")
	serviceAddr := flag.String("service-addr", "", "service URL")
	configPath := flag.String("conf", "", "Config path")
	flag.Parse()
	
	var config Config
	
	data, err := os.ReadFile(*configPath)
	
	if err != nil {
		fmt.Printf("Error: %s \n", err.Error())
		os.Exit(1)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		fmt.Printf("Error: failed to parse config file \n")
		os.Exit(1)
	}
	url, err := url.Parse(*serviceAddr)
	if err != nil {
		fmt.Printf("Error: %s \n", err.Error())
		os.Exit(1)
	}
	
	proxy := &httputil.ReverseProxy{
		Rewrite: func(req *httputil.ProxyRequest) {
			req.SetURL(url)
		},
		Transport: NewFirewall(config),
	}

	fmt.Printf("Server listening on %s ...\n", *addr)
	http.ListenAndServe(*addr, proxy)
}

func forbidden() *http.Response {
	body := "Forbidden"
	return &http.Response{
		StatusCode:    http.StatusForbidden,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Header:        http.Header{"Content-Type": []string{"text/plain"}},
		Body:          io.NopCloser(strings.NewReader(body)),
		ContentLength: int64(len(body)),
	}
}

func allowRequest(header http.Header, body []byte, rule Rule) bool {
	// forbidden headers are "Name: Value" pairs
	for _, h := range rule.ForbiddenHeaders {
		parts := strings.SplitN(h, ": ", 2)
		if len(parts) == 2 && header.Get(parts[0]) == parts[1] {
			return false
		}
	}

	for _, k := range rule.RequiredHeaders {
		if header.Get(k) == "" {
			return false
		}
	}

	// forbidden user agents are regex patterns matched against User-Agent value
	ua := header.Get("User-Agent")
	for _, pattern := range rule.ForbiddenUserAgents {
		if regexp.MustCompile(pattern).MatchString(ua) {
			return false
		}
	}

	if rule.MaxRequestLengthBytes > 0 {
		length, err := strconv.Atoi(header.Get("Content-Length"))
		if err == nil && length > rule.MaxRequestLengthBytes {
			return false
		}
	}

	for _, pattern := range rule.ForbiddenRequestRe {
		if regexp.MustCompile(pattern).Match(body) {
			return false
		}
	}

	return true
}

func allowResponse(statusCode int, body []byte, rule Rule) bool {
	if slices.Contains(rule.ForbiddenResponseCodes, statusCode) {
		return false
	}

	if rule.MaxResponseLengthBytes > 0 && len(body) > rule.MaxResponseLengthBytes {
		return false
	}

	for _, pattern := range rule.ForbiddenResponseRe {
		if regexp.MustCompile(pattern).Match(body) {
			return false
		}
	}

	return true
}

