package main

import (
	"database/sql"
	"time"
)

type StockItem struct {
	StockID           int64
	ProductName       string
	QuantityRemaining float64
	ExpirationDate    time.Time
}

func GetAllProducts(db *sql.DB) ([]Product, error) {
	rows, err := db.Query("SELECT id,barcode,name,kcal_per_100,protein_per_100,fat_per_100,carbs_per_100,unit_type,category FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product

	for rows.Next() {
		var p Product
		err := rows.Scan(&p.ID, &p.Barcode, &p.Name, &p.KcalPer100, &p.ProteinPer100, &p.FatPer100, &p.CarbsPer100, &p.UnitType, &p.Category)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}
func InsertProduct(db *sql.DB, p Product) (int64, error) {
	var id int64
	err := db.QueryRow(
		"INSERT INTO products (barcode, name, kcal_per_100, protein_per_100, fat_per_100,carbs_per_100, unit_type, category ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8 ) RETURNING id",
		p.Barcode, p.Name, p.KcalPer100, p.ProteinPer100, p.FatPer100, p.CarbsPer100, p.UnitType, p.Category,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func InsertStockEntry(db *sql.DB, s StockEntry) (int64, error) {
	var id int64
	err := db.QueryRow("INSERT INTO stock_entries (product_id,quantity_remaining,expiration_date) VALUES ( $1, $2, $3 ) RETURNING id",
		s.ProductID, s.QuantityRemaining, s.ExpirationDate,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func GetAllStock(db *sql.DB) ([]StockItem, error) {
	rows, err := db.Query("SELECT stock_entries.id, products.name, stock_entries.quantity_remaining, stock_entries.expiration_date\n FROM stock_entries\n JOIN products ON stock_entries.product_id = products.id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stockItem []StockItem

	for rows.Next() {
		var s StockItem
		err := rows.Scan(&s.StockID, &s.ProductName, &s.QuantityRemaining, &s.ExpirationDate)
		if err != nil {
			return nil, err
		}
		stockItem = append(stockItem, s)
	}
	return stockItem, nil
}
func ConsumeStock(db *sql.DB, stockID int64, amount float64) (float64, error) {
	var remaining float64
	err := db.QueryRow(
		"UPDATE stock_entries SET quantity_remaining = quantity_remaining - $1 WHERE id = $2 RETURNING quantity_remaining",
		amount, stockID,
	).Scan(&remaining)
	if err != nil {
		return 0, err
	}
	return remaining, nil
}
func SetStockQuantity(db *sql.DB, stockID int64, newQuantity float64) (float64, error) {
	var remaining float64
	err := db.QueryRow(
		"UPDATE stock_entries SET quantity_remaining = $1 WHERE id = $2 RETURNING quantity_remaining",
		newQuantity, stockID,
	).Scan(&remaining)
	if err != nil {
		return 0, err
	}
	return remaining, nil
}
func MarkStockFinished(db *sql.DB, stockID int64) error {
	_, err := db.Exec("UPDATE stock_entries SET quantity_remaining = 0 WHERE id = $1", stockID)
	if err != nil {
		return err
	}
	return nil
}
