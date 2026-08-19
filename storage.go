package main

import (
	"database/sql"
)

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
func InsertProduct(db *sql.DB, p Product) (int64, error) { //функция которая берет соединение с БД и продукт а возвращает id и статус ошибки
	var id int64        // создаем переменную id
	err := db.QueryRow( //вставляем продукт в бд
		"INSERT INTO products (barcode, name, kcal_per_100, protein_per_100, fat_per_100,carbs_per_100, unit_type, category ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8 ) RETURNING id",
		p.Barcode, p.Name, p.KcalPer100, p.ProteinPer100, p.FatPer100, p.CarbsPer100, p.UnitType, p.Category,
	).Scan(&id) //сканируем id который она она отдала
	if err != nil {
		return 0, err
	} //проверяем на ошибку
	return id, nil //возвращаем значения
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
