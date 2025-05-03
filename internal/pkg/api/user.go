package api

import (
	"net/http"
	"slices"
	"strconv"

	"github.com/wilgnert/mnbe/internal/pkg/middleware"
	"github.com/wilgnert/mnbe/internal/pkg/services"
)

const (
	GetAllUsers  = "GetAllUsers"
	RegisterUser = "RegisterUser"
	GetUserByID  = "GetUserByID"
	Login        = "Login"
)

func (c *Config) handlerLogin(w http.ResponseWriter, r *http.Request) {
	// provided by middleware
	email := r.Context().Value(middleware.Key(middleware.UserEmail)).(string)
	password := r.Context().Value(middleware.Key(middleware.UserPassword)).(string)

	status, payload, err := services.LoginWithEmail(email, password, c.secret, c.db)

	if err != nil {
		respondWithError(w, status, err.Error())
	} else {
		respondWithJSON(w, status, payload)
	}
}

func (c *Config) handlerGetUserByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("user_id")

	status, payload, err := services.GetUserById(id, c.db)

	if err != nil {
		respondWithError(w, status, err.Error())
	} else {
		respondWithJSON(w, status, payload)
	}
}

func (c *Config) handlerGetAllUsers(w http.ResponseWriter, r *http.Request) {
	var page = 1
	var limit = 20
	var offset = 0
	var sort = "newest"
	sortTypes := []string{"newest", "oldest", "lexical_asc", "lexical_desc"}
	if r.URL.Query().Has("page") {
		p, err := strconv.Atoi(r.URL.Query().Get("page"))
		if err == nil && p > 0 {
			page = p
		}
	}
	if r.URL.Query().Has("limit") {
		l, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err == nil && l > 0 {
			if l > 100 {
				l = 100
			}
			limit = l
		}
	}
	if r.URL.Query().Has("offset") {
		o, err := strconv.Atoi(r.URL.Query().Get("offset"))
		if err == nil && o > 0 {
			if o > limit {
				o = limit
			}
			offset = o
		}
	}
	if r.URL.Query().Has("sort") {
		s := r.URL.Query().Get("sort")
		if slices.Contains(sortTypes, s) {
			sort = s
		}
	}

	status, payload, err := services.GetAllUsers(page, limit, offset, sort, c.db)

	if err != nil {
		respondWithError(w, status, err.Error())
	} else {
		respondWithJSON(w, status, payload)
	}
}

func (c *Config) handlerRegisterUser(w http.ResponseWriter, r *http.Request) {
	email := r.Context().Value(middleware.Key(middleware.UserEmail)).(string)
	password := r.Context().Value(middleware.Key(middleware.UserPassword)).(string)

	status, payload, err := services.CreateUserWithEmail(email, password, c.db)

	if err != nil {
		respondWithError(w, status, err.Error())
	} else {
		respondWithJSON(w, status, payload)
	}

}
