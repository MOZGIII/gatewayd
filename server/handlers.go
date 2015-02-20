package server

import (
	"net/http"
	"strconv"

	"gatewayd/api"
	"gatewayd/tunnel"

	"gatewayd/pkg/encoder"
	"github.com/go-martini/martini"
)

// NewPublicHander creates a handler with public-faced services
// There are WebSocket tunnel and public APIs
func NewPublicHander() http.Handler {
	handler := http.NewServeMux()

	handler.HandleFunc("/tunnel", tunnel.Handler)

	return handler
}

// NewServiceHander created a handler with private-faced services
// These are mostly internal APIs
func NewServiceHander() http.Handler {
	m := buildMartini()

	// Sessions
	m.Post("/sessions", api.CreateSession)
	m.Get("/sessions/:token", api.SessionByToken)
	m.Delete("/sessions/:token", api.DeleteSession)

	return m
}

func buildMartini() *martini.ClassicMartini {
	r := martini.NewRouter()
	m := martini.New()
	m.Use(martini.Logger())
	m.Use(martini.Recovery())
	m.MapTo(r, (*martini.Routes)(nil))

	// Map encoder conditionally based on pretty print request param
	m.Use(func(c martini.Context, w http.ResponseWriter, r *http.Request) {
		pretty, _ := strconv.ParseBool(r.URL.Query().Get("pretty"))
		c.MapTo(encoder.JSONEncoder{PrettyPrint: pretty}, (*encoder.Encoder)(nil))
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	})

	m.Action(r.Handle)
	return &martini.ClassicMartini{Martini: m, Router: r}
}
