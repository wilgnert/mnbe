package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/wilgnert/mnbe/internal/pkg/lib"
	"github.com/wilgnert/mnbe/internal/pkg/middleware"
)

type Config struct {
	db lib.Database
	Handlers map[string]http.Handler
	secret string
}

func (c *Config) SetDB(db lib.Database) {
	c.db = db
}

func (c *Config) SetSecret(s string) {
	c.secret = s
}

func (c *Config) RegisterHandlers() {
	if c.Handlers == nil {
		c.Handlers = make(map[string]http.Handler)
	}
	c.Handlers[GetAllUsers] = http.HandlerFunc(c.handlerGetAllUsers)
	c.Handlers[RegisterUser] = middleware.CreateStack(
		middleware.ExtractEmail,
		middleware.ExtractPassword,
	)(http.HandlerFunc(c.handlerRegisterUser))
	c.Handlers[GetUserByID] = http.HandlerFunc(c.handlerGetUserByID)
	c.Handlers[Login] = middleware.CreateStack(
		middleware.ExtractEmail,
		middleware.ExtractPassword,
	)(http.HandlerFunc(c.handlerLogin))

}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) error {
	res, err := json.Marshal(payload)
	if err != nil {
		return err;
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	w.Write(res)
	return nil
}

func respondWithError(w http.ResponseWriter, code int, msg string) error {
	return respondWithJSON(w, code, map[string]string{"error": msg})
}

func RespondOK(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	io.WriteString(w, "OK")
}

func RespondNoContent(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}