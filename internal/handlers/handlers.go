package handlers

import (
	"net/http"

	"github.com/ironedr/gochat/internal/application"
)

type Handler struct {
	app *application.Application
}

func New(app *application.Application) *Handler {
	return &Handler{app}
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	td := h.NewTemplateData(r)
	h.Render(w, http.StatusOK, "home", td)
}
