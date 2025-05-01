package server

import (
	"crypto/tls"
	"net/http"
)

type Server struct {
	http.Server
}

func NewServer(addr string, router http.Handler, tlsConfig *tls.Config) *Server {
	server := Server{
		http.Server{
			Addr: addr,
			Handler: router,
			TLSConfig: tlsConfig,
		},
	}
	return &server
}