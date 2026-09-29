package main

import "github.com/shravan26/serverhost/registry"

type Server struct {
	registry *registry.Registry
}

func NewServer() *Server {
	return &Server{
		registry: registry.New(),
	}
}
