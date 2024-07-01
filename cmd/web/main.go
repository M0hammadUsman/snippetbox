package main

import (
	"context"
	"crypto/tls"
	"flag"
	"github.com/M0hammadUsman/snippetbox/internal/models"
	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lmittmann/tint"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"
)

// Add a snippets field to the application struct. This will allow us to
// make the SnippetModel object available to our handlers.
type application struct {
	snippets       *models.SnippetModel
	users          *models.UserModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	// Configs
	type config struct {
		addr      string
		staticDir string
		dsn       string
	}
	var cfg config
	flag.StringVar(&cfg.addr, "addr", ":8080", "Http Network Address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.StringVar(&cfg.dsn, "dsn", "postgres://web:uusmann3344@localhost:5432/snippetbox?sslmode=disable", "Data Source Name (DSN)")
	flag.Parse()
	// Loggers
	slog.SetDefault(slog.New(tint.NewHandler(os.Stderr, &tint.Options{AddSource: true})))
	// Database
	dbPool, err := openDB(cfg.dsn)
	if err != nil {
		log.Fatal("Unable to create connection pool\n", err.Error())
	}
	defer dbPool.Close()
	tmplCache, err := newTemplateCache()
	if err != nil {
		log.Fatal("Unable to cache the templates\n", err.Error())
	}
	// Configuring session manager
	sessionManager := scs.New()
	sessionManager.Store = pgxstore.New(dbPool)
	sessionManager.Lifetime = 12 * time.Hour
	// Initialize a models.SnippetModel instance and add it to the application dependencies.
	app := &application{
		snippets:       &models.SnippetModel{DBPool: dbPool},
		users:          &models.UserModel{DBPool: dbPool},
		templateCache:  tmplCache,
		formDecoder:    form.NewDecoder(),
		sessionManager: sessionManager,
	}
	// TLS config to restrict the curves
	tlsConfig := &tls.Config{CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256}}
	srv := &http.Server{
		Addr:      cfg.addr,
		Handler:   app.routes(),
		TLSConfig: tlsConfig,
		// Add Idle, Read and Write timeouts to the server.
		IdleTimeout:  100 * time.Minute,
		ReadTimeout:  1000 * time.Second,
		WriteTimeout: 1000 * time.Second,
	}
	slog.Info("Starting server on", "address", cfg.addr)
	slog.Info("Static dir path set to", "path", cfg.staticDir)
	slog.Error(srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem").Error())
}

func openDB(dsn string) (*pgxpool.Pool, error) {
	// to fine-tune the pool config first create config using pgx.parseConfig and then use pgx.NewWithConfig
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	if err = pool.Ping(context.Background()); err != nil {
		return nil, err
	}
	return pool, nil
}
