package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

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
	host := r.URL.Query().Get("host")
	hostParts := strings.Split(host, ",")
	port := r.URL.Query().Get("port")
	portParts := strings.Split(port, ",")
	if len(hostParts) != len(portParts) {
		http.Error(w, "host and port count mismatch", http.StatusBadRequest)
	}
	var serviceInstances []registry.ServiceInstance
	for i := range hostParts {
		serviceInstances = append(serviceInstances, registry.ServiceInstance{
			Host: hostParts[i],
			Port: portParts[i],
		})
	}

	service := r.URL.Query().Get("service")

	s.registry.Register(service, serviceInstances)
	w.WriteHeader(http.StatusCreated)
}
func (s *Server) getServiceInstances(w http.ResponseWriter, r *http.Request) {
	service := r.URL.Query().Get("service")
	instances, avl := s.registry.Next(service)
	if avl != true {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Service not available"))
	}
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
