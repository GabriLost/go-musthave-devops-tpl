package server

import "github.com/go-chi/chi/v5"

func Router() chi.Router {
	router := chi.NewRouter()
	router.Get("/", GetAllHandler)
	router.Get("/value/{typ}/{name}", GetMetricHandler)
	router.Post("/update/{typ}/{name}/{value}", PostMetricHandler)
	router.Post("/update/{typ}/", NotFoundHandler)
	router.Post("/update/*", NotImplementedHandler)
	router.Post("/update/", PostJsonMetricsHandler)
	//router.Post("/value/", PostJsonMetricsHandler)
	return router
}
