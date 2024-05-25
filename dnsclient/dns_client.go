package dnsclient

import (
	"github.com/miekg/dns"
	"net/http"
	"time"
)

type DnsClient interface {
	Exchange(req *dns.Msg) (resp *dns.Msg, ttt time.Duration, err error)
}

type dnsClient struct {
	name    string
	address string

	httpClient *http.Client
}
