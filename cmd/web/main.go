package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"

	"emmasela-snippetbox/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresDB struct {
	pool *pgxpool.Pool
}

// Add a snippets field to the application struct. This will allow us to
// make the SnippetModel object available to our handlers.

type application struct {
	errorLog *log.Logger
	infoLog *log.Logger
	snippets *models.SnippetModel
}



func main() {


	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "postgresql://web:emmasela@localhost:5432/snippetbox", "PostgreSQL data source name")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()


	// Initialize a models.SnippetModel instance and add it to the application
	// dependencies.
	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
		snippets: &models.SnippetModel{DB: db.pool},
	}

	server := &http.Server{
		Addr: *addr,
		ErrorLog: errorLog,
		Handler: app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = server.ListenAndServe()
	errorLog.Fatal(err)


}

func openDB(dsn string) (*postgresDB, error) {
	db, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(context.Background()); err != nil {
		return nil, err
	}

	return &postgresDB{pool: db}, nil
}

func (db *postgresDB) Close() {
	db.pool.Close()
}

func (db *postgresDB) Ping(ctx context.Context) error {
	return db.pool.Ping(ctx)
}
