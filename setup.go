package alias

import (
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"

	"github.com/caddyserver/caddy"
)

// func init() {
// 	caddy.RegisterPlugin("alias", caddy.Plugin{
// 		ServerType: "dns",
// 		Action:     setup,
// 	})
// }

// init registers this plugin.
func init() { plugin.Register("example", setup) }

// setup is the function that gets called when the config parser see the token "example". Setup is responsible
// for parsing any extra options the example plugin may have. The first token this function sees is "example".
func setup(c *caddy.Controller) error {
	c.Next()
	if c.NextArg() {
		return plugin.Error("alias", c.ArgErr())
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return Alias{Next: next}
	})

	return nil
}
