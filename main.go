package main

import (
	"database/sql"
	"log"
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
	http.HandleFunc("/products/list", app.getProductsHandler)
	http.HandleFunc("/stock", app.addStockHandler)
	http.HandleFunc("/stock/list", app.getStockHandler)
	http.HandleFunc("/stock/consume", app.consumeStockHandler)
	http.HandleFunc("/stock/set", app.setStockHandler)
	http.HandleFunc("/stock/finish", app.markStockHandler)
	http.HandleFunc("/stock/expiring", app.expiringStockHandler)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
