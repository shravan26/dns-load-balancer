package main

import (
	"github.com/shravan26/serverhost/registry"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func loggingMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf(
			"MiniCloud hit: method=%s uri=%s host=%s",
			r.Method,
			r.RequestURI,
			r.Host,
		)
		next.ServeHTTP(w, r)
	})
}
func main() {
	target, err := registry.RouteTo()
	if err != nil {
		log.Fatal(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(target)
	server := loggingMiddleWare(proxy)
	log.Println("MiniCloud server running on :8080")

	err = http.ListenAndServe("0.0.0.0:8080", server)
	if err != nil {
		log.Fatal(err)
	}
}
