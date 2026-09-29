package registry

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func reconstructUrl(r *http.Request, sI []ServiceInstance) {
	scheme := "http"
	if r.Header.Get("X-Forwarded-Proto") != "" {
		scheme = r.Header.Get("X-Forwarded-Proto")
	} else if r.TLS != nil {
		scheme = "https"
	}
	r.URL.Scheme = scheme
	r.URL.Host = sI[0].Host + ":" + sI[0].Port

	fmt.Printf("Read request as %s\n", r.URL)
}
func (s *Server) RouteTo() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/")
		service := strings.Split(path, "/")[0]
		serviceInstances := s.registry.Get(service)
		if len(serviceInstances) == 0 {
			http.Error(w, "Service not registered", http.StatusBadGateway)
			return
		}

		reconstructUrl(r, serviceInstances)
		target, err := url.Parse(r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		proxy := httputil.NewSingleHostReverseProxy(target)
		proxy.ServeHTTP(w, r)
	})
}
