// package custom_loadbalance is a plugin for rewriting responses to do "load balancing"
package weighted_loadbalance

import (
	"context"

	"github.com/coredns/coredns/plugin"

	"github.com/miekg/dns"
)

type Percentage struct {
	Next    plugin.Handler
	Configs LoadBalancesFile
}

func (p Percentage) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	wrr := &PercentageResponseWriter{Configs: p.Configs, ResponseWriter: w}
	return plugin.NextOrFailure(p.Name(), p.Next, ctx, wrr, r)
}

// Name implements the Handler interface.
func (p Percentage) Name() string { return "weighted_loadbalance(percentage)" }
