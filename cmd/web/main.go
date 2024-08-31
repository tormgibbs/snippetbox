package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"emmasela-snippetbox/internal/models"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresDB struct {
	pool *pgxpool.Pool
}

// Add a snippets field to the application struct. This will allow us to
// make the SnippetModel object available to our handlers.

type application struct {
	errorLog 				*log.Logger
	infoLog 				*log.Logger
	snippets 				*models.SnippetModel
	formDecoder 		*form.Decoder
	sessionManager 	*scs.SessionManager
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

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = pgxstore.New(db.pool)
	sessionManager.Lifetime = 12 * time.Hour

	// Setting the session cookie to secure to ensure that it is sent when HTTPS connection is used
	sessionManager.Cookie.Secure = true


	// Initialize a models.SnippetModel instance and add it to the application
	// dependencies.
	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
		snippets: &models.SnippetModel{DB: db.pool},
		formDecoder: formDecoder,
		sessionManager: sessionManager,
	}

	server := &http.Server{
		Addr: *addr,
		ErrorLog: errorLog,
		Handler: app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = server.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
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
