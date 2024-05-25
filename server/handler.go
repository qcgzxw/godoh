package server

import (
	"errors"
	"fmt"
	"github.com/miekg/dns"
	"github.com/qcgzxw/godoh/dnsclient"
	"time"
)

type handler struct {
	dohClient dnsclient.DnsClient
}

func (this *handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	var (
		questionName string
		qClass       uint16
		qType        uint16
	)
	if err := func(w dns.ResponseWriter, r *dns.Msg) (err error) {
		if len(r.Question) != 1 {
			err = errors.New("number of questions is not 1")
			return
		}
		questionName = r.Question[0].Name
		qClass = r.Question[0].Qclass
		qType = r.Question[0].Qclass
		var (
			resp *dns.Msg
			ttt  time.Duration
		)
		if resp, ttt, err = this.dohClient.Exchange(r); err != nil {
			return
		}
		fmt.Printf("[%d]question: %s, class: %s, type: %s cost: %s\n", r.Id, questionName, dns.ClassToString[qClass], dns.TypeToString[qType], ttt.String())
		_ = w.WriteMsg(resp)
		return
	}(w, r); err != nil {
		var questionClass, questionType string
		if qClass > 0 {
			questionClass = dns.ClassToString[qClass]
		}
		if qType > 0 {
			questionType = dns.TypeToString[qType]
		}
		fmt.Printf(
			"%s - - [%s] \"%s %s %s\"\n",
			w.RemoteAddr().String(),
			time.Now().Format("02/Jan/2006:15:04:05 -0700"),
			questionName,
			questionClass,
			questionType,
		)
		reply := new(dns.Msg)
		reply.Rcode = dns.RcodeFormatError
		_ = w.WriteMsg(reply)
		return
	}
}

func NewHandler(dohClient dnsclient.DnsClient) dns.Handler {
	return &handler{
		dohClient: dohClient,
	}
}
