package main

import (
	"github.com/wilgnert/mnbe/internal/pkg/api"
	"github.com/wilgnert/mnbe/internal/pkg/lib"
	"github.com/wilgnert/mnbe/internal/pkg/routes"
	"github.com/wilgnert/mnbe/internal/pkg/server"
)

func main() {
	db := lib.StubDatabase{}
	cfg := api.Config{}
	cfg.SetDB(&db)
	cfg.RegisterHandlers()

	router := routes.LoadRoutes(&cfg)
	s := server.NewServer(":8080", router, nil)
	s.ListenAndServe()
}