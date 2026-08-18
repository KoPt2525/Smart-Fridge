package main

import (
	"database/sql"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "postgres://localhost/postgres?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	app := &App{db: db}
	http.HandleFunc("/products", app.addProductHandler)
	http.ListenAndServe("localhost:8080", nil)
}
