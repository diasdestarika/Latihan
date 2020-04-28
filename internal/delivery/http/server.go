package http

import (
	"net/http"

	"github.com/rs/cors"
	"github.com/vilbert/go-skeleton/pkg/grace"
)

// UserHandler ...
type UserHandler interface {
	UserHandler(w http.ResponseWriter, r *http.Request)
}

// Server ...
type Server struct {
	server *http.Server
	User   UserHandler
}

// Serve is serving HTTP gracefully on port x ...
func (s *Server) Serve(port string) error {
	handler := cors.AllowAll().Handler(s.Handler()) //biar bisa di askes dari yang lain, dan lebih aman
	return grace.Serve(port, handler)
}
