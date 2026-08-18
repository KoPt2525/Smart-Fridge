package main

import "database/sql"

type OFFResponse struct {
	Code    string
	Status  int
	Product OFFProduct
}

type OFFProduct struct {
	ProductName string `json:"product_name"`
	Brands      string
	Nutriments  OFFNutriments
}

type OFFNutriments struct {
	Carbohydrates100g float64 `json:"carbohydrates_100g"`
	Proteins100g      float64 `json:"proteins_100g"`
	Fat100g           float64 `json:"fat_100g"`
	EnergyKcal100g    float64 `json:"energy-kcal_100g"`
	EnergyKj100g      float64 `json:"energy-kj_100g"`
}

func mapToProduct(off OFFResponse, barcode string) Product {
	if off.Product.Nutriments.EnergyKcal100g == 0 {
		off.Product.Nutriments.EnergyKcal100g = off.Product.Nutriments.EnergyKj100g / 4.184
	}
	p := Product{
		Barcode:       barcode,
		Name:          off.Product.ProductName,
		KcalPer100:    off.Product.Nutriments.EnergyKcal100g,
		ProteinPer100: off.Product.Nutriments.Proteins100g,
		FatPer100:     off.Product.Nutriments.Fat100g,
		CarbsPer100:   off.Product.Nutriments.Carbohydrates100g,
		UnitType:      UnitWeight,
	}
	return p
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
