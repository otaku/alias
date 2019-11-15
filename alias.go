package alias

import (
	"context"

	"github.com/coredns/coredns/plugin"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/miekg/dns"
)

var log = clog.NewWithPlugin("alias")

// Rewrite is plugin to rewrite requests internally before being handled.
type Alias struct {
	Next plugin.Handler
}

// ServeDNS implements the plugin.Handler interface.
func (a Alias) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	mw := NewResponseModifier(w)
	return plugin.NextOrFailure(a.Name(), a.Next, ctx, mw, r)
}

// Name implements the Handler interface.
func (a Alias) Name() string { return "alias" }

type ResponseModifier struct {
	dns.ResponseWriter
}

func (r *ResponseModifier) WriteMsg(res *dns.Msg) error {
	// Find CNAMEs

	found := false
	for i := 0; i < len(res.Answer); {
		answer := res.Answer[i]

		if answer.Header().Rrtype != dns.TypeCNAME {
			i++
			continue
		}

		cname := answer.(*dns.CNAME)

		// Find A record
		found = false
		for j, rr := range res.Answer {
			if rr.Header().Rrtype == dns.TypeA && rr.Header().Name == cname.Target {
				res.Answer[j].Header().Name = answer.Header().Name
				found = true
			}
		}

		// Remove CNAME
		if found {
			res.Answer = append(res.Answer[:i], res.Answer[i+1:]...)
			continue
		}

		i++
	}

	return r.ResponseWriter.WriteMsg(res)
}

func NewResponseModifier(w dns.ResponseWriter) *ResponseModifier {
	return &ResponseModifier{
		ResponseWriter: w,
	}
}
