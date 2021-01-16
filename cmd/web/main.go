package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net/http"
	"os"
	"raissov/snippetbox/pkg/models/postgresql"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	snippets *postgresql.SnippetModel
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	dbpool, err := pgxpool.Connect(context.Background(), "host=localhost port=5432 user=postgres password=postgres dbname=snippetbox sslmode=disable")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		snippets: &postgresql.SnippetModel{DB: dbpool},
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
