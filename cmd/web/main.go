package main

import (
	"context"
	"flag"
	"github.com/M0hammadUsman/snippetbox/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lmittmann/tint"
	"log"
	"log/slog"
	"net/http"
	"os"
)

// Add a snippets field to the application struct. This will allow us to
// make the SnippetModel object available to our handlers.
type application struct {
	snippets *models.SnippetModel
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
	// Initialize a models.SnippetModel instance and add it to the application dependencies.
	app := &application{
		snippets: &models.SnippetModel{DBPool: dbPool},
	}
	defer dbPool.Close()
	slog.Info("Starting server on", "address", cfg.addr)
	slog.Info("Static dir path set to", "path", cfg.staticDir)
	slog.Error(http.ListenAndServe(cfg.addr, app.routes()).Error())
}

func openDB(dsn string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	if err = pool.Ping(context.Background()); err != nil {
		return nil, err
	}
	return pool, nil
}
