package handlers

import (
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

func (h *Handler) SetupRoutes() *mux.Router {
	r := mux.NewRouter()
	r.Use(h.ConfigureHeaders)
	r.Use(h.CSRF)
	r.NotFoundHandler = http.HandlerFunc(h.NotFound)
	r.HandleFunc("/", h.Home).Methods(http.MethodGet)

	return r
}

// ConfigureHeaders sets sensible secure default headers such as Content-Security-Poicy, X-Content-Type-Options, X-Frame-Options and X-XSS-Protections
func (h *Handler) ConfigureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com https://cdnjs.cloudflare.com/ajax/; font-src fonts.gstatic.com  https://cdnjs.cloudflare.com/ajax/; script-src 'self' unpkg.com https://cdn.jsdelivr.net/npm/")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) CSRF(next http.Handler) http.Handler {
	CSRF := csrf.Protect(h.app.CSRFKey, csrf.HttpOnly(true), csrf.SameSite(csrf.SameSiteStrictMode), csrf.Path("/"), csrf.Secure(true))
	return CSRF(next)
}
