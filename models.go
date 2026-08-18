package main

import "time"

type UnitType string

const (
	UnitWeight UnitType = "weight"
	UnitCount  UnitType = "count"
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
	Category      string
}

type StockEntry struct {
	ID                int64
	ProductID         int64
	QuantityRemaining float64
	ExpirationDate    time.Time
	AddedAt           time.Time
}
