package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Proxy struct {
	url *url.URL
}

func New(target string) (*Proxy, error) {
	targetUrl, err := url.Parse(target)
	if err != nil {
		return nil, err
	}
	return &Proxy{targetUrl}, nil
}

func (p *Proxy) serveReverseProxy(w http.ResponseWriter, r *http.Request) {

	proxy := httputil.NewSingleHostReverseProxy(p.url)

	r.URL.Host = p.url.Host
	r.URL.Scheme = p.url.Scheme
	r.Host = p.url.Host
	proxy.ServeHTTP(w, r)
}

func (p *Proxy) HandleProxyRequest(w http.ResponseWriter, r *http.Request) {
	p.serveReverseProxy(w, r)
}
