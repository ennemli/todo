package middlewares

import (
	"context"
	"net/http"
	"time"

	"github.com/ennemli/todo/user/internal/errors"
	"github.com/go-chi/render"
)

func SetTimeOut(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			done := make(chan bool)
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()
			go func() {
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
				done <- true
			}()
			select {
			case <-done:
				return
			case <-ctx.Done():
				renderErrorRequestTimeout(w, r)

			}
		})
	}
}

func renderErrorRequestTimeout(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusRequestTimeout)
	render.JSON(w, r, errors.ResponseRequestTimeout)
}
