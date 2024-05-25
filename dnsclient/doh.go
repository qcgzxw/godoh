package dnsclient

import (
	"encoding/base64"
	"github.com/miekg/dns"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type dohDnsClient dnsClient

// Exchange implement rfc8484 https://www.rfc-editor.org/rfc/rfc8484#section-6
func (this *dohDnsClient) Exchange(reqDns *dns.Msg) (respDns *dns.Msg, ttt time.Duration, err error) {

	var req *http.Request
	if req, err = http.NewRequest("GET", this.address, nil); err != nil {
		log.Fatal(err)
		return
	}

	req.Header.Add("accept", "application/dns-message")
	req.Header.Add("content-type", "application/dns-message")

	q := req.URL.Query()

	var (
		buf    []byte
		begin  = time.Now()
		origID = reqDns.Id
	)
	reqDns.Id = 0
	if buf, err = reqDns.Pack(); err != nil {
		log.Fatal(err)
		return
	}
	q.Add("dns", base64.RawURLEncoding.EncodeToString(buf))

	req.URL.RawQuery = q.Encode()

	var resp *http.Response
	if resp, err = this.httpClient.Do(req); err != nil {
		log.Fatal(err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatal("Upstream DNS server error.")
	}

	defer resp.Body.Close()
	var body []byte
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		log.Fatal(err)
		return
	}
	respDns = new(dns.Msg)
	if err = respDns.Unpack(body); err != nil {
		log.Fatal(err)
		return
	}
	respDns.Id = origID
	ttt = time.Since(begin)
	return
}

func NewDoHDnsClient(name, path string) DnsClient {
	return &dohDnsClient{
		name:    name,
		address: path,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}
func NewJDoHDnsClient(name, path string) DnsClient {
	return &jDohDnsClient{
		name:    name,
		address: path,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}
