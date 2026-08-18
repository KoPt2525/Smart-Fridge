package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "postgres://localhost/postgres?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	off, err := GetProductByBarcode("4610003254766")
	if err != nil {
		panic(err)
	}

	p := mapToProduct(off, "4610003254766")

	id, err := InsertProduct(db, p)
	if err != nil {
		panic(err)
	}

	fmt.Printf("вставлено с id=%d, продукт: %+v\n", id, p)
}
