package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/shravan26/serverhost/registry"
)

type Server struct {
	registry *registry.Registry
}

func checkEmptyPath(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "" {
			http.Error(w, "No empty path value functions", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello there, welcome to this server")
}

func (s *Server) register(w http.ResponseWriter, r *http.Request) {
	instance := registry.ServiceInstance{
		Host: r.URL.Query().Get("host"),
		Port: r.URL.Query().Get("port"),
	}
	service := r.URL.Query().Get("service")

	s.registry.Register(service, instance)
	w.WriteHeader(http.StatusCreated)
}
func (s *Server) getServiceInstances(w http.ResponseWriter, r *http.Request) {
	service := r.URL.Query().Get("service")
	instances := s.registry.Get(service)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(instances)
}
func main() {
	server := &Server{
		registry: registry.New(),
	}
	result := http.NewServeMux()
	result.HandleFunc("/hello", handler)
	result.HandleFunc("/register", server.register)
	result.HandleFunc("/services", server.getServiceInstances)
	serveHandler := checkEmptyPath(result)
	err := http.ListenAndServe(":9000", serveHandler)
	if err != nil {
		log.Fatal(err)
	}
}
