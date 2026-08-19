package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

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

func GetProductByBarcode(barcode string) (OFFResponse, error) {
	url := fmt.Sprintf("https://world.openfoodfacts.org/api/v2/product/%s.json?fields=product_name,brands,nutriments,status,code", barcode)
	resp, err := http.Get(url)
	if err != nil {
		return OFFResponse{}, err
	}
	defer resp.Body.Close()
	var off OFFResponse
	err = json.NewDecoder(resp.Body).Decode(&off)
	if err != nil {

		return OFFResponse{}, err
	}
	if off.Status == 0 {
		return OFFResponse{}, fmt.Errorf("product not found: %s", barcode)
	}
	return off, nil
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
