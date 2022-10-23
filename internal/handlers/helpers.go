package handlers

import (
	"net/http"
	"runtime/debug"
)

// ServerError logs an internal error message and stack trace and sends a status 500 response to the user
func (h *Handler) ServerError(w http.ResponseWriter, err error) {
	h.app.Logger.Printf("%s\n%s\n", err.Error(), debug.Stack())
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// ClientError is used for errors caused by the incoming request and sends a client specific status code to the user
func (h *Handler) ClientError(w http.ResponseWriter, status int, err error) {
	h.app.Logger.Println(err)
	http.Error(w, http.StatusText(status), status)
}

func (h *Handler) NotFound(w http.ResponseWriter, r *http.Request) {
	td := h.NewTemplateData(r)
	h.Render(w, http.StatusNotFound, "notfound", td)
}
