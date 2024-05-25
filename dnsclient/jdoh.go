package dnsclient

import (
	"encoding/json"
	"fmt"
	"github.com/miekg/dns"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type jDohDnsClient dohDnsClient

func (this *jDohDnsClient) Exchange(reqDns *dns.Msg) (respDns *dns.Msg, ttt time.Duration, err error) {
	var (
		begin  = time.Now()
		origID = reqDns.Id
	)
	reqDns.Id = 0
	jsonResponse := this.Lookup(reqDns.Question[0].Name, reqDns.Question[0].Qtype)
	respDns = DohJsonResponseToDnsMsg(*jsonResponse)
	respDns.Id = origID

	ttt = time.Since(begin)
	return
}

// Lookup implement DoH json-api
func (this *jDohDnsClient) Lookup(name string, rType uint16) (jsonResponse *DohJsonResponse) {

	req, err := http.NewRequest("GET", this.address, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("accept", "application/dns-json")

	q := req.URL.Query()
	q.Add("name", name)
	q.Add("type", strconv.Itoa(int(rType)))
	q.Add("cd", "false") // unset DNSSEC
	q.Add("do", "false") // ignore disable validation
	req.URL.RawQuery = q.Encode()

	resp, err := this.httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatal("Upstream DNS server error.")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("DNS RESPONSE BODY:\n%s\n", body)

	jsonResponse = new(DohJsonResponse)
	err = json.Unmarshal(body, jsonResponse)
	if err != nil {
		log.Fatal(err)
	}

	return
}
