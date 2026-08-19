package main

import (
	"time"
)

type UnitType string

const (
	UnitWeight UnitType = "weight"
	UnitCount  UnitType = "count"
)

type Category string

const (
	CategoryMilkAndEggs Category = "milk products and eggs"
	CategoryMeat        Category = "meat and poultry"
	CategoryFish        Category = "fish and seafood"
	CategoryVegetables  Category = "vegetables"
	CategoryFruits      Category = "fruits and berries"
	CategoryBakery      Category = "bakery"
	CategoryGrain       Category = "pasta and grains"
	CategoryCanned      Category = "canned goods"
	CategoryFrozen      Category = "frozen foods"
	CategorySauces      Category = "sauces and spices"
	CategorySweets      Category = "sweets and snacks"
	CategoryDrinks      Category = "drinks"
	CategoryOils        Category = "oils and fats"
	CategoryOther       Category = "other"
)

type Product struct {
	ID            int64
	Barcode       string
	Name          string
	KcalPer100    float64
	ProteinPer100 float64
	FatPer100     float64
	CarbsPer100   float64
	UnitType      UnitType
	Category      Category
}

type StockEntry struct {
	ID                int64
	ProductID         int64
	QuantityRemaining float64
	ExpirationDate    time.Time
	AddedAt           time.Time
}
