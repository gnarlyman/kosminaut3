package anon

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/proxy"
)

const defaultUserAgent = "Mozilla/5.0 (compatible)"

type Config struct {
	ProxyURL  string
	UserAgent string
	Timeout   time.Duration
}

func NewClient(cfg Config) (*http.Client, error) {
	base := &net.Dialer{Timeout: 5 * time.Second, KeepAlive: 30 * time.Second}

	dialContext := base.DialContext
	if cfg.ProxyURL != "" {
		u, err := url.Parse(cfg.ProxyURL)
		if err != nil {
			return nil, fmt.Errorf("parse proxy url: %w", err)
		}
		if u.Scheme != "socks5" {
			return nil, fmt.Errorf("unsupported proxy scheme %q: only socks5 is supported", u.Scheme)
		}
		d, err := proxy.FromURL(u, base)
		if err != nil {
			return nil, fmt.Errorf("build socks5 dialer: %w", err)
		}
		cd, ok := d.(proxy.ContextDialer)
		if !ok {
			return nil, fmt.Errorf("socks5 dialer does not support context")
		}
		dialContext = cd.DialContext
	}

	transport := &http.Transport{
		DialContext:           dialContext,
		ForceAttemptHTTP2:     true,
		TLSHandshakeTimeout:   5 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		IdleConnTimeout:       90 * time.Second,
		MaxIdleConns:          10,
	}

	ua := cfg.UserAgent
	if ua == "" {
		ua = defaultUserAgent
	}

	return &http.Client{
		Transport: &scrubTransport{inner: transport, userAgent: ua},
		Timeout:   cfg.Timeout,
	}, nil
}

type scrubTransport struct {
	inner     http.RoundTripper
	userAgent string
}

func (s *scrubTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	clone := req.Clone(reqContext(req))
	if clone.Header == nil {
		clone.Header = make(http.Header)
	}
	clone.Header.Set("User-Agent", s.userAgent)
	clone.Header.Del("Referer")
	clone.Header.Del("From")
	clone.Header.Del("X-Forwarded-For")
	clone.Header.Del("X-Real-IP")
	return s.inner.RoundTrip(clone)
}

func reqContext(req *http.Request) context.Context {
	if ctx := req.Context(); ctx != nil {
		return ctx
	}
	return context.Background()
}
