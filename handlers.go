package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
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
type AddStockRequest struct {
	ProductID         int64     `json:"product_id"`
	QuantityRemaining float64   `json:"quantity_remaining"`
	ExpirationDate    time.Time `json:"expiration_date"`
}
type AddStockResponse struct {
	ID    int64      `json:"id"`
	Stock StockEntry `json:"stock"`
}

type ConsumeStockRequest struct {
	StockID int64   `json:"stock_id"`
	Amount  float64 `json:"amount"`
}

type SetStockRequest struct {
	StockID     int64   `json:"stock_id"`
	NewQuantity float64 `json:"new_quantity"`
}

type MarkFinishedRequest struct {
	StockID int64 `json:"stock_id"`
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
func (a *App) getProductsHandler(w http.ResponseWriter, r *http.Request) {
	products, err := GetAllProducts(a.db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(products)
}

func (a *App) addStockHandler(w http.ResponseWriter, r *http.Request) {
	var req AddStockRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s := StockEntry{
		ProductID:         req.ProductID,
		QuantityRemaining: req.QuantityRemaining,
		ExpirationDate:    req.ExpirationDate,
	}
	id, err := InsertStockEntry(a.db, s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(AddStockResponse{ID: id, Stock: s})
}
func (a *App) getStockHandler(w http.ResponseWriter, r *http.Request) {
	stockItem, err := GetAllStock(a.db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(stockItem)
}

func (a *App) consumeStockHandler(w http.ResponseWriter, r *http.Request) {
	var req ConsumeStockRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	remaining, err := ConsumeStock(a.db, req.StockID, req.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]float64{"remaining": remaining})
}
func (a *App) setStockHandler(w http.ResponseWriter, r *http.Request) {
	var req SetStockRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	remaining, err := SetStockQuantity(a.db, req.StockID, req.NewQuantity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]float64{"remaining": remaining})
}
func (a *App) markStockHandler(w http.ResponseWriter, r *http.Request) {
	var req MarkFinishedRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = MarkStockFinished(a.db, req.StockID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func (a *App) expiringStockHandler(w http.ResponseWriter, r *http.Request) {
	items, err := GetExpiringStock(a.db, 10)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(items)
}
