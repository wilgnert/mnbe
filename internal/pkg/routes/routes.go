package routes

import "net/http"

func LoadRoutes() *http.ServeMux {
	router := http.NewServeMux()

	v1 := http.NewServeMux()
	v1.Handle("/v1/", http.StripPrefix("/v1", router))

	return router
}