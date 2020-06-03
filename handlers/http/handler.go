package http

import (
	"github.com/gorilla/mux"
)

// Handler provides interface for handlers
type Handler interface {
	LoadRoutes(*mux.Router)
}
