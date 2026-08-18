package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type AddProductRequest struct {
	Barcode string `json:"barcode"`
}
type AddProductResponse struct {
	ID      int64   `json:"id"`
	Product Product `json:"product"`
}

type App struct {
	db *sql.DB
}

func (a *App) addProductHandler(w http.ResponseWriter, r *http.Request) {
	var req AddProductRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	off, err := GetProductByBarcode(req.Barcode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	p := mapToProduct(off, req.Barcode)

	id, err := InsertProduct(a.db, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(AddProductResponse{ID: id, Product: p})

}
