package options

import (
	"github.com/miekg/dns"
)

type Options struct {
	DnsServers             []dns.Server
	DnsServerConfiguration DnsServerConfiguration
	DnsCacheConfiguration  DnsCacheConfiguration

	UpstreamDnsServers             []dns.Client
	UpstreamDnsServerConfiguration UpstreamDnsServerConfiguration
}

type DnsServerConfiguration struct {
	// The number of requests per second allowed per client. Setting it to 0 means no limit.
	RateLimit int
	// Add the EDNS Client Subnet option (ECS) to upstream requests
	// and log the values sent by the clients in the query log.
	EDNS bool
	// Set DNSSEC flag in the outcoming DNS queries and check the result (DNSSEC-enabled resolver is required).
	DNSSEC bool
	// Disable resolving of IPv6 addresses
	DisableIPV6 bool
}
type DnsCacheConfiguration struct {
	// DNS cache size (in bytes). To disable caching, leave empty.
	CacheSize int64
}

type UpstreamDnsServerConfiguration struct {
	// balance mode: 0 - random, 1 - round-robin, 2 - sequential
	BalanceMode uint8
}

func DefaultOptions() Options {
	return Options{
		DnsServers:                     nil,
		DnsServerConfiguration:         DnsServerConfiguration{},
		DnsCacheConfiguration:          DnsCacheConfiguration{},
		UpstreamDnsServers:             nil,
		UpstreamDnsServerConfiguration: UpstreamDnsServerConfiguration{},
	}
}

func NewOptionsFromConfig() Options {
	return DefaultOptions()
}
