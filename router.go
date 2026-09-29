package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/shravan26/serverhost/registry"
)

func reconstructUrl(r *http.Request, sI []registry.ServiceInstance, path string) {
	scheme := "http"
	if r.Header.Get("X-Forwarded-Proto") != "" {
		scheme = r.Header.Get("X-Forwarded-Proto")
	} else if r.TLS != nil {
		scheme = "https"
	}
	r.URL.Scheme = scheme
	r.URL.Host = sI[0].Host + ":" + sI[0].Port
	r.URL.Path = path
	fmt.Printf("Read request as %s\n", r.URL)
}
func (s *Server) RouteTo() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/")
		parts := strings.Split(path, "/")
		service := parts[0]
		fmt.Println(service)
		serviceInstances := s.registry.Get(service)
		if len(serviceInstances) == 0 {
			http.Error(w, "Service not registered", http.StatusBadGateway)
			return
		}
		forwardPath := strings.Join(parts[1:], "/")
		reconstructUrl(r, serviceInstances, forwardPath)
		target, err := url.Parse(
			r.URL.Scheme + "://" + r.URL.Host,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		proxy := httputil.NewSingleHostReverseProxy(target)
		proxy.ServeHTTP(w, r)
	})
}
