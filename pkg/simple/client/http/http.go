package http

import (
	"net"
	"net/http"
	"time"
)

var Client *http.Client

func Init(timeout time.Duration) {
	Client = NewHttpClient(timeout)
}

func NewHttpClient(timeout time.Duration) *http.Client {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		DisableKeepAlives:     false,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          500,
		MaxIdleConnsPerHost:   500,
		MaxConnsPerHost:       2000,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	// http2.ConfigureTransport(transport)

	return &http.Client{
		Timeout:   timeout,
		Transport: transport,
	}
}
