package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/gormstore"
	"github.com/alexedwards/scs/v2"
	"github.com/ironedr/gochat/internal/application"
	"github.com/ironedr/gochat/internal/handlers"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	var port int
	var dsn, baseURL, env, csrfKey string
	flag.IntVar(&port, "port", 9000, "server port")
	flag.StringVar(&dsn, "dsn", "postgres://myuser:test123@localhost:5432/gochat", "postgesql DSN")
	flag.StringVar(&baseURL, "base-url", "localhost/", "base-url used for static files")
	flag.StringVar(&csrfKey, "csrf-key", "my5ecret", "CSRF secret key")
	flag.Parse()
	app := &application.Application{}
	app.Logger = log.New(os.Stdout, "GoChat", log.Ldate|log.Ltime)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		app.Logger.Fatal(err)
	}
	templateCache, err := handlers.NewTemplateCache()
	if err != nil {
		app.Logger.Fatal(err)
	}

	sessionManager := scs.New()
	sessionManager.Store, err = gormstore.New(db)
	if err != nil {
		app.Logger.Fatal(err)
	}
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true
	sessionManager.Cookie.SameSite = http.SameSiteStrictMode

	app.BaseURL = baseURL
	app.DSN = dsn
	app.Env = env
	app.CSRFKey = []byte(csrfKey)
	app.SessionManager = sessionManager
	app.TemplateCache = templateCache
	handler := handlers.New(app)
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.Port),
		Handler:      handler.SetupRoutes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	app.Logger.Fatal(srv.ListenAndServeTLS("tls/cert.pem", "tls/key.pem"))
}
