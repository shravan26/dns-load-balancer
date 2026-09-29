package main

import (
	"log"
	"net/http"

	"github.com/shravan26/serverhost/registry"
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
	server := NewServer()
	server.registry.Register("api", registry.ServiceInstance{
		Host: "localhost",
		Port: "9000",
	})
	target := server.RouteTo()
	actedRequest := loggingMiddleWare(target)
	log.Println("MiniCloud server running on :8080")

	err := http.ListenAndServe("0.0.0.0:8080", actedRequest)
	if err != nil {
		log.Fatal(err)
	}
}
