// package weighted_loadbalance shuffles A, AAAA and MX records.
package weighted_loadbalance

import (
	"fmt"

	"github.com/coredns/coredns/plugin/pkg/log"
	"github.com/miekg/dns"
)

type PercentageResponseWriter struct {
	dns.ResponseWriter
	Configs LoadBalancesFile
}

// WriteMsg implements the dns.ResponseWriter interface.
func (r *PercentageResponseWriter) WriteMsg(res *dns.Msg) error {
	if res.Rcode != dns.RcodeSuccess {
		return r.ResponseWriter.WriteMsg(res)
	}

	if res.Question[0].Qtype == dns.TypeAXFR || res.Question[0].Qtype == dns.TypeIXFR {
		return r.ResponseWriter.WriteMsg(res)
	}

	res.Answer = r.roundRobin(res.Answer, res.Question)
	// res.Ns = roundRobin(res.Ns)
	// res.Extra = roundRobin(res.Extra)

	return r.ResponseWriter.WriteMsg(res)
}

func (r *PercentageResponseWriter) roundRobin(in []dns.RR, question []dns.Question) []dns.RR {
	cname := []dns.RR{}
	address := []dns.RR{}
	mx := []dns.RR{}
	rest := []dns.RR{}

	out := []dns.RR{}

	for _, r := range in {
		switch r.Header().Rrtype {
		case dns.TypeCNAME:
			cname = append(cname, r)
		case dns.TypeA:
			address = append(address, r)
		case dns.TypeMX:
			mx = append(mx, r)
		default:
			rest = append(rest, r)
		}
	}

	if len(address) > 0 {
		var subDomain string
		particles := Strings.Split(question[0].Name, ".")
		if len(particles) > 2 {
			subDomain = particles[0]
		} else {
			subDomain = "root"
		}
		wieghtedSelection(r.Configs.RecordTypes.A[subDomain])
	}
	out = append(cname, rest...)
	out = append(out, address...)
	out = append(out, mx...)
	return out
}

func wieghtedSelection(contents []Content) {
	fmt.Println(contents)
}

// Write implements the dns.ResponseWriter interface.
func (r *PercentageResponseWriter) Write(buf []byte) (int, error) {
	// Should we pack and unpack here to fiddle with the packet... Not likely.
	log.Warning("RoundRobin called with Write: not shuffling records")
	n, err := r.ResponseWriter.Write(buf)
	return n, err
}
