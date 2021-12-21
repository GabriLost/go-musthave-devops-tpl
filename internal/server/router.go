package server

import "github.com/go-chi/chi/v5"

func Router() chi.Router {
	router := chi.NewRouter()
	router.Get("/", GetAllHandler)
	router.Get("/value/{typ}/{name}", ValueMetricHandler)
	router.Post("/value/", JSONValueHandler)
	router.Post("/update/{typ}/{name}/{value}", UpdateMetricHandler)
	router.Post("/update/{typ}/{name}", BadRequestHandler)
	router.Post("/update/{typ}/", NotFoundHandler)
	router.Post("/update/", JSONUpdateMetricsHandler)
	router.Post("/update/*", NotImplementedHandler)
	return router
}
