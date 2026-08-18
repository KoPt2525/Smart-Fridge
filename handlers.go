package main

import (
	"encoding/json"
	"net/http"
)

type AddProductRequest struct {
	Barcode string `json:"barcode"`
}

func addProductHandler(w http.ResponseWriter, r *http.Request) {
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
	json.NewEncoder(w).Encode(p)

}
