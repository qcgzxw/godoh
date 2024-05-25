package dnsclient

import (
	"github.com/miekg/dns"
	"net"
)

const (
	RegularDns int64 = iota + 1
	DnsOverUDP
	DnsOverTCP
	DnsOverTLS
	DnsOverHTTPS
	DnsOverHTTPSWithForcedH3
	DnsOverQUIC
)

// DoH JSON schema using by Cloudflare and Google
// Cloudflare’s DNS over HTTPS endpoint also supports JSON format for querying DNS data. For lack of an agreed upon JSON schema for DNS over HTTPS in the Internet Engineering Task Force (IETF), Cloudflare has chosen to follow the same schema as Google’s DNS over HTTPS resolver.
type DohJsonRequest struct {
	Name string `json:"name"`
	Type uint16 `json:"type"`
	DO   bool   `json:"do,omitempty"`
	CD   bool   `json:"cd,omitempty"`
}

type DohJsonResponse struct {
	Status           int           `json:"Status"`
	TC               bool          `json:"TC"`
	RD               bool          `json:"RD"`
	RA               bool          `json:"RA"`
	AD               bool          `json:"AD"`
	CD               bool          `json:"CD"`
	Question         []Question    `json:"Question"`
	Answer           []Answer      `json:"Answer"`
	Authority        []Answer      `json:"Authority,omitempty"`
	Additional       []interface{} `json:"Additional,omitempty"`
	EdnsClientSubnet string        `json:"edns_client_subnet,omitempty"`
	Comment          string        `json:"Comment,omitempty"`
}

type Question struct {
	Name string `json:"name"`
	Type uint16 `json:"type"`
}

type Answer struct {
	Name string `json:"name"`
	Type uint16 `json:"type"`
	TTL  int    `json:"TTL"`
	Data string `json:"data"`
}

// 将 DoH JSON 响应转换为 dns.Msg
func DohJsonResponseToDnsMsg(dohResp DohJsonResponse) *dns.Msg {
	msg := new(dns.Msg)
	msg.Response = true
	msg.Authoritative = true
	msg.RecursionAvailable = dohResp.RA
	msg.AuthenticatedData = dohResp.AD
	msg.CheckingDisabled = dohResp.CD
	msg.Rcode = dohResp.Status

	// 转换 Question
	for _, q := range dohResp.Question {
		msg.Question = append(msg.Question, dns.Question{
			Name:   dns.Fqdn(q.Name),
			Qtype:  q.Type,
			Qclass: dns.ClassINET,
		})
	}

	// 转换 Answer
	for _, ans := range dohResp.Answer {
		msg.Answer = append(msg.Answer, convertAnswerToRR(ans))
	}

	// 转换 Authority
	for _, auth := range dohResp.Authority {
		msg.Ns = append(msg.Ns, convertAnswerToRR(auth))
	}

	// 转换 Additional (这里假设 Additional 是空的)
	// 你需要根据实际情况处理 Additional 字段

	return msg
}

// 辅助函数：将 Answer 转换为 dns.RR
func convertAnswerToRR(ans Answer) dns.RR {
	hdr := dns.RR_Header{
		Name:   dns.Fqdn(ans.Name),
		Rrtype: ans.Type,
		Class:  dns.ClassINET,
		Ttl:    uint32(ans.TTL),
	}

	switch ans.Type {
	case dns.TypeA:
		return &dns.A{
			Hdr: hdr,
			A:   parseARecord(ans.Data),
		}
	case dns.TypeAAAA:
		return &dns.AAAA{
			Hdr:  hdr,
			AAAA: parseAAAARecord(ans.Data),
		}
	case dns.TypeCNAME:
		return &dns.CNAME{
			Hdr:    hdr,
			Target: dns.Fqdn(ans.Data),
		}
	case dns.TypeMX:
		// MX 记录需要优先级，这里假设优先级为 10
		return &dns.MX{
			Hdr:        hdr,
			Preference: 10,
			Mx:         dns.Fqdn(ans.Data),
		}
	case dns.TypeTXT:
		return &dns.TXT{
			Hdr: hdr,
			Txt: []string{ans.Data},
		}
	// 添加其他类型的记录处理
	default:
		return nil
	}
}

// 辅助函数：解析 A 记录的 IP 地址
func parseARecord(data string) net.IP {
	return net.ParseIP(data).To4()
}

// 辅助函数：解析 AAAA 记录的 IP 地址
func parseAAAARecord(data string) net.IP {
	return net.ParseIP(data).To16()
}
