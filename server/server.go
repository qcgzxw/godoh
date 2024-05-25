package server

import (
	"github.com/miekg/dns"
	"github.com/qcgzxw/godoh/dnsclient"
	"log"
	"strings"
)

type Server interface {
	ListenAndServe() error
}

type dnsClient struct {
	dohClient dnsclient.DnsClient

	udpServers []*dns.Server
	tcpServers []*dns.Server
}

func (c *dnsClient) ListenAndServe() error {
	results := make(chan error, len(c.udpServers)+len(c.tcpServers))
	for _, srv := range append(c.udpServers, c.tcpServers...) {
		go func(srv *dns.Server) {
			err := srv.ListenAndServe()
			if err != nil {
				log.Println(err)
			}
			results <- err
		}(srv)
	}

	for i := 0; i < cap(results); i++ {
		err := <-results
		if err != nil {
			return err
		}
	}
	close(results)

	return nil
}

func NewServer(conf *Config) Server {
	c := &dnsClient{}

	c.dohClient = dnsclient.NewDoHDnsClient("", conf.Upstream)

	dnsHandler := NewHandler(c.dohClient)
	for _, addr := range conf.Listen {
		println("godoh is running at: " + addr)
		if strings.Contains(addr, "tcp://") {
			addr = strings.TrimPrefix(addr, "tcp://")
			c.tcpServers = append(c.tcpServers, &dns.Server{
				Addr:    addr,
				Net:     "tcp",
				Handler: dnsHandler,
			})
		} else if strings.Contains(addr, "udp://") {
			addr = strings.TrimPrefix(addr, "udp://")
			c.udpServers = append(c.udpServers, &dns.Server{
				Addr:    addr,
				Net:     "udp",
				Handler: dnsHandler,
				UDPSize: dns.DefaultMsgSize,
			})
		} else {
			c.tcpServers = append(c.tcpServers, &dns.Server{
				Addr:    addr,
				Net:     "tcp",
				Handler: dnsHandler,
			})
			c.udpServers = append(c.udpServers, &dns.Server{
				Addr:    addr,
				Net:     "udp",
				Handler: dnsHandler,
				UDPSize: dns.DefaultMsgSize,
			})
		}
	}
	return c
}
