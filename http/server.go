package http

import (
	"net"
	"net/http"
)

// Server ...
type Server struct {
	ln      net.Listener
	Handler *Handler
	Addr    string
}

// Open starts listening on a socket and serves the HTTP server
func (s *Server) Open() error {
	// opens a socket, bind it to an address, and start listening
	ln, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}
	s.ln = ln

	// start HTTP server
	go func() { http.Serve(s.ln, s.Handler) }()

	return nil
}

// Close closes the socket
func (s *Server) Close() error {
	if s.ln != nil {
		s.ln.Close()
	}
	return nil
}
