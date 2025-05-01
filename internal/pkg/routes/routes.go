package routes

import (
	"net/http"

	"github.com/wilgnert/mnbe/internal/pkg/api"
)

func LoadRoutes(cfg *api.Config) *http.ServeMux {
	user := http.NewServeMux()
	user.Handle("GET /users/", cfg.Handlers[api.GetAllUsers])
	user.Handle("POST /users/", cfg.Handlers[api.RegisterUser])

	v0 := http.NewServeMux()
	v0.Handle("/v0/", http.StripPrefix("/v0", user))
	v0.Handle("/", http.HandlerFunc(api.RespondOK))

	return v0
}