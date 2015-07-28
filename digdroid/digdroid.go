package digdroid

//Package digdroid is Android and (maybe) iOS binding for running DNS test using pulse.
//We don;t use pulse directly cause bind doesnt work with types such as interface{} and uint16
//gomobile bind github.com/turbobytes/pulse/digdroid

import (
	"github.com/miekg/dns"
	"github.com/turbobytes/pulse/utils"
	"log"
)

type DNSResult struct {
	Err    string
	Output string
	Rtt    string
}

//We use RunDNS as proxy instead of using pulse directly because gomobile bind
//can't make bindings for types that use interface{} or uint16
func RunDNS(host, target, qtypestr string, norec bool) *DNSResult {
	var qtype uint16
	switch qtypestr {
	case "A":
		qtype = 1
	case "AAAA":
		qtype = 28
	case "NS":
		qtype = 2
	case "CNAME":
		qtype = 5
	}
	req := &pulse.DNSRequest{
		Host:        host,
		QType:       qtype,
		Targets:     []string{target},
		NoRecursion: norec,
	}
	result := pulse.DNSImpl(req)
	res := &DNSResult{}
	if result.Err != "" {
		res.Err = result.Err
	} else if len(result.Results) < 1 {
		res.Err = "No results returned"
	} else if result.Results[0].Err != "" {
		res.Err = result.Results[0].Err
		res.Rtt = result.Results[0].RttStr
	} else {
		msg := &dns.Msg{}
		msg.Unpack(result.Results[0].Raw)
		res.Output = msg.String()
		log.Println(res.Output)
		res.Rtt = result.Results[0].RttStr
	}
	return res
}
