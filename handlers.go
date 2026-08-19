package main

import ( //пакеты для работы с бд энкодерами и сетью
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

type AddProductRequest struct { //структра которая принимает баркод из json
	Barcode string `json:"barcode"`
}
type AddProductResponse struct { //структура которая принимает id и продукты
	ID      int64   `json:"id"`
	Product Product `json:"product"`
}

type App struct { //структура которая держит ссылку на БД?
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

func (a *App) addProductHandler(w http.ResponseWriter, r *http.Request) { //метод который работает с БД и отдает запросы к
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
