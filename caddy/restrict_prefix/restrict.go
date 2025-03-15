package restrictprefix

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"go.uber.org/zap"
)

func init() {
	caddy.RegisterModule(RestrictPrefix{})
}

type RestrictPrefix struct {
	Prefix string `json:"prefix,omitempty"`
	logger *zap.Logger
}

// CaddyModule returns the CaddyModule infomation
func (RestrictPrefix) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.restrict_prefix",
		New: func() caddy.Module { return new(RestrictPrefix) },
	}
}

// Provision a Zap logger to RestrictPrefix
func (p *RestrictPrefix) Provision(ctx caddy.Context) error {
	p.logger = ctx.Logger(p)
	return nil
}

// Validate the prefix from the modules configuration, setting the
// default prefix "." if necessary
func (p *RestrictPrefix) Validate() error {
	if p.Prefix == "" {
		p.Prefix = "."
	}
	return nil
}

func (p RestrictPrefix) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	for _, part := range strings.Split(r.URL.Path, "/") {
		if strings.HasPrefix(part, p.Prefix) {
			http.Error(w, "Not Found", http.StatusNotFound)
			if p.logger != nil {
				p.logger.Debug(fmt.Sprintf("restrict prefix: %q in %s\n", part, r.URL.Path))
			}
			return nil
		}
	}
	return next.ServeHTTP(w, r)
}

var (
	_ caddy.Provisioner           = (*RestrictPrefix)(nil)
	_ caddy.Validator             = (*RestrictPrefix)(nil)
	_ caddyhttp.MiddlewareHandler = (*RestrictPrefix)(nil)
)
