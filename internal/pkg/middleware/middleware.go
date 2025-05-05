package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/wilgnert/mnbe/internal/pkg/auth"
)

const (
	UserEmail = "user.email"
	UserPassword = "user.password"
	UserID = "user.id"
	RequestAuthorization = "request.Authorization"
	RequestApiKey = "request.api-key"
	RequestBearer = "request.bearer"
)

type Key string

type Middleware func(http.Handler) http.Handler

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int 
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func CreateStack(xs ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(xs) - 1; i >= 0; i-- {
			x := xs[i]
			next = x(next)
		}
		return next
	}
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := &wrappedWriter {
			ResponseWriter: w,
			statusCode: http.StatusOK,
		}
		next.ServeHTTP(wrapped, r)
		log.Println(wrapped.statusCode, r.Method, r.URL.Path, time.Since(start))
	})
}

func ExtractFieldFromBodyMiddleware(fieldName string, ExtractedFieldKey Key) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Only process POST, PUT, or PATCH requests with JSON content type.
			if r.Method != http.MethodPost && r.Method != http.MethodPut && r.Method != http.MethodPatch {
				next.ServeHTTP(w, r)
				return
			}
			if r.Header.Get("Content-Type") != "application/json" {
				http.Error(w, "Content-Type must be application/json", http.StatusBadRequest)
				return
			}

			// Read the entire body into a buffer.
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Failed to read request body", http.StatusBadRequest)
				return
			}

			// Restore the request body for subsequent handlers.
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			var body map[string]string
			err = json.Unmarshal(bodyBytes, &body)
			if err != nil {
				http.Error(w, "Failed to decode JSON body", http.StatusBadRequest)
				return
			}

			if value, ok := body[fieldName]; ok {
				// Create a new context with the extracted value.
				ctx := context.WithValue(r.Context(), ExtractedFieldKey, value)
				// Serve the next handler with the new context.
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				http.Error(w, fmt.Sprintf("Field '%s' not found in JSON body", fieldName), http.StatusBadRequest)
				return
			}
		})
	}
}

func ExtractFieldFromHeaderMiddleware(fieldName string, extractedFieldKey Key) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract the value from the header.
			value := r.Header.Get(fieldName)

			// Check if the header field exists.
			if value == "" {
				http.Error(w, fmt.Sprintf("Header field '%s' not found", fieldName), http.StatusBadRequest)
				return // Important:  Return after sending the error.
			}

			// Create a new context with the extracted value.
			ctx := context.WithValue(r.Context(), extractedFieldKey, value)

			// Serve the next handler with the new context.
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func ExtractAPIKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		value, err := auth.GetAPIKey(r.Header)

		// Check if the header field exists.
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return // Important:  Return after sending the error.
		}
		// Create a new context with the extracted value.
		ctx := context.WithValue(r.Context(), Key(RequestApiKey), value)

		// Serve the next handler with the new context.
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func ExtractBearer(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		v, err := auth.GetBearerToken(r.Header)

		// Check if the header field exists.
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return // Important:  Return after sending the error.
		}

		// Create a new context with the extracted value.
		ctx := context.WithValue(r.Context(), Key(RequestBearer), v)

		// Serve the next handler with the new context.
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

var ExtractEmail = ExtractFieldFromBodyMiddleware("email", UserEmail)
var ExtractPassword = ExtractFieldFromBodyMiddleware("password", UserPassword)
var ExtractUserID = ExtractFieldFromBodyMiddleware("user_id", UserID)