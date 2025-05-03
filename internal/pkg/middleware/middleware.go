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
)

const (
	UserEmail = "user.email"
	UserPassword = "User.password"
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

func ExtractFieldMiddleware(fieldName string, ExtractedFieldKey Key) func(http.Handler) http.Handler {
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

var ExtractEmail = ExtractFieldMiddleware("email", UserEmail)
var ExtractPassword = ExtractFieldMiddleware("password", UserPassword)