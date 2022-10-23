package application

import (
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
)

// Application is used to pass common core dependencies to all our subpackages
type Application struct {
	TemplateCache  map[string]*template.Template
	Logger         *log.Logger
	SessionManager *scs.SessionManager
	Env            string
	DSN            string
	BaseURL        string
	CSRFKey        []byte
	Port           int
}
