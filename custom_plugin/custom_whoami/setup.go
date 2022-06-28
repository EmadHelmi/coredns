package custom_whoami

import (
	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
)

func init() { plugin.Register("custom_whoami", setup) }

func setup(c *caddy.Controller) error {
	c.Next() // 'whoami'
	if c.NextArg() {
		return plugin.Error("custom_whoami", c.ArgErr())
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return Whoami{Next: next}
	})

	return nil
}
