package main

import (
	"github.com/qcgzxw/godoh/server"
)

func main() {
	// simple test DoH server
	// >bash "nslookup -port=1234 -type=A example.com"
	s := server.NewServer(&server.Config{
		Listen:   []string{"udp://0.0.0.0:1234"},
		Upstream: "https://unfiltered.adguard-dns.com/dns-query",
	})
	_ = s.ListenAndServe()
}
