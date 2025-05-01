package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/wilgnert/mnbe/internal/pkg/lib"
	"github.com/wilgnert/mnbe/internal/pkg/middleware"
)

const (
	GetAllUsers = "GetAllUsers"
	RegisterUser = "RegisterUser"
)

type Config struct {
	db lib.Database
	Handlers map[string]http.Handler
}

func (c *Config) SetDB(db lib.Database) {
	c.db = db
}

func (c *Config) RegisterHandlers() {
	if c.Handlers == nil {
		c.Handlers = make(map[string]http.Handler)
	}
	c.Handlers[GetAllUsers] = http.HandlerFunc(c.handlerGetAllUsers)
	c.Handlers[RegisterUser] = middleware.CreateStack(
		middleware.ExtractEmail,
		middleware.ExtractHashedPassword,
	)(http.HandlerFunc(c.handlerRegisterUser))
}

func (c *Config) handlerGetAllUsers(w http.ResponseWriter, r *http.Request) {
	allUsers, err := c.db.RetrieveUsers()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	respondWithJSON(w, http.StatusOK, allUsers)
}
func (c *Config) handlerRegisterUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Context())
	email := r.Context().Value(middleware.Key(middleware.UserEmail)).(string)
	hashed_password := r.Context().Value(middleware.Key(middleware.UserHashedPassword)).(string)

	user, err := c.db.CreateUser(lib.CreateUserParams{
		Email: email,
		HashedPassword: hashed_password,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, user)

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