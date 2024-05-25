package server

import (
	"testing"
)

func TestNewDnsClient(t *testing.T) {
	conf := &Config{
		Listen:   []string{"udp://0.0.0.0:1234"},
		Upstream: "https://unfiltered.adguard-dns.com/dns-query",
	}
	c := NewServer(conf)
	if err := c.ListenAndServe(); err != nil {
		t.Error(err)
	}
}
