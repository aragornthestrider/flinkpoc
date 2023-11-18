package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Route struct {
	Name        string
	Pattern     string
	Method      string
	HandlerFunc http.HandlerFunc
}

func (s *Server) NewRouter() *mux.Router {
	routes := []Route{
		{
			"test",
			"/test",
			http.MethodGet,
			http.HandlerFunc(s.GetTest),
		},
	}

	supportRoutes := []Route{
		{
			"liveness",
			"/liveness",
			http.MethodGet,
			http.HandlerFunc(s.CheckLive),
		},
		{
			"readiness",
			"/readiness",
			http.MethodGet,
			http.HandlerFunc(s.CheckReady),
		},
		{
			"Metrics - prometheus metrics for monitoring",
			"/metrics",
			http.MethodGet,
			promhttp.Handler().ServeHTTP,
		},
	}

	router := mux.NewRouter()

	for _, route := range routes {
		mainRouter := router.PathPrefix("").Subrouter().StrictSlash(true)
		mainRouter.Name(route.Name).Path(route.Pattern).Methods(route.Method).Handler(route.HandlerFunc)
		mainRouter.Use(prometheusMiddleware)
	}

	for _, route := range supportRoutes {
		supportRouter := router.PathPrefix("").Subrouter().StrictSlash(true)
		supportRouter.Name(route.Name).Path(route.Pattern).Methods(route.Method).Handler(route.HandlerFunc)
	}

	return router
}
