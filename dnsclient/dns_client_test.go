package dnsclient

import (
	"github.com/miekg/dns"
	"log"
	"testing"
)

func newGoogleDNS() DnsClient {
	return NewDoHDnsClient("Google", "https://dns.google/dns-query")
}

func newCloudFlareDNS() DnsClient {
	return NewDoHDnsClient("CloudFlare", "https://cloudflare-dns.com/dns-query")
}

func newAdGuardDNS() DnsClient {
	return NewDoHDnsClient("AdGuard", "https://dns.adguard.com/dns-query")
}

func newDiegoDNS() DnsClient {
	return NewDoHDnsClient("Diego", "https://public.diego.run/dns-query")
}

func newCloudFlareJsonDNS() DnsClient {
	return NewJDoHDnsClient("CloudFlare", "https://cloudflare-dns.com/dns-query")
}
func newGoogleJsonDNS() DnsClient {
	return NewJDoHDnsClient("CloudFlare", "https://dns.google/resolve")
}

// 生成新的 DNS 查询消息
func newQueryDNSMsg(name string, qType uint16) *dns.Msg {
	return &dns.Msg{
		MsgHdr: dns.MsgHdr{
			RecursionDesired: true,
		},
		Question: []dns.Question{
			{
				Name:   dns.Fqdn(name),
				Qtype:  qType,
				Qclass: dns.ClassINET,
			},
		},
	}
}

// Helper function to test DoH clients
func testDoHDnsClient(t *testing.T, client DnsClient, name string, qType uint16) {
	query := newQueryDNSMsg(name, qType)
	response, ttt, err := client.Exchange(query)
	if err != nil {
		t.Fatalf("Failed to query %s: %v", name, err)
	}

	if len(response.Answer) == 0 {
		t.Fatalf("No answers received for %s", name)
	}
	log.Printf("cost: %s", ttt.String())
	for _, answer := range response.Answer {
		log.Printf("Answer: %v", answer)
	}
}
func TestDoHDnsClient(t *testing.T) {
	// testDoHDnsClient(t, newGoogleDNS(), "example.com", dns.TypeA) // blocked by GFW
	testDoHDnsClient(t, newCloudFlareDNS(), "example.com.", dns.TypeA)
	testDoHDnsClient(t, newAdGuardDNS(), "example.com.", dns.TypeA)
	testDoHDnsClient(t, newDiegoDNS(), "example.com.", dns.TypeA)
}
func TestJDoHDnsClient(t *testing.T) {
	// testDoHDnsClient(t, newGoogleJsonDNS(), "example.com", dns.TypeA) // blocked by GFW
	testDoHDnsClient(t, newCloudFlareJsonDNS(), "example.com", dns.TypeA)
}
