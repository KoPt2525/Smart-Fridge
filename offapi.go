package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type OFFResponse struct { //создаем структуру которая содержит код статус и другую структуру
	Code    string
	Status  int
	Product OFFProduct
}

type OFFProduct struct { //структура которая берет имя с json бренд и структуру нутриентов
	ProductName string `json:"product_name"`
	Brands      string
	Nutriments  OFFNutriments
}

type OFFNutriments struct { //структура содержащая энергетическую ценность и КБЖУ
	Carbohydrates100g float64 `json:"carbohydrates_100g"`
	Proteins100g      float64 `json:"proteins_100g"`
	Fat100g           float64 `json:"fat_100g"`
	EnergyKcal100g    float64 `json:"energy-kcal_100g"`
	EnergyKj100g      float64 `json:"energy-kj_100g"`
}

func GetProductByBarcode(barcode string) (OFFResponse, error) { //функция которая принимает баркоде а отдает структуру
	url := fmt.Sprintf("https://world.openfoodfacts.org/api/v2/product/%s.json?fields=product_name,brands,nutriments,status,code", barcode) //составляем запрос
	resp, err := http.Get(url)                                                                                                              //получаем ответ и проверяем ответ сервера на ошибку
	if err != nil {
		return OFFResponse{}, err
	} // проверяем чтобы ошибки не было, в случае ошибки возвращаем пустую структуру
	defer resp.Body.Close()                       //закрываем подключение к апи(в конце функции)
	var off OFFResponse                           //создаем структуру offresponce под именем off
	err = json.NewDecoder(resp.Body).Decode(&off) //декодируем структуру по ссылке на off
	if err != nil {

		return OFFResponse{}, err
	} //проверка что удалось ли распарсить json
	if off.Status == 0 {
		return OFFResponse{}, fmt.Errorf("product not found: %s", barcode)
	} //проверяем что продукт нашелся
	return off, nil //возвращаем структуру off
}

func mapToProduct(off OFFResponse, barcode string) Product { //функция которая принимает структуру продукта и баркод а возвращает структуру продукта
	if off.Product.Nutriments.EnergyKcal100g == 0 {
		off.Product.Nutriments.EnergyKcal100g = off.Product.Nutriments.EnergyKj100g / 4.184
	} //перевод джоулей в каллории
	p := Product{ //инициализируем новый продукт
		Barcode:       barcode,
		Name:          off.Product.ProductName,
		KcalPer100:    off.Product.Nutriments.EnergyKcal100g,
		ProteinPer100: off.Product.Nutriments.Proteins100g,
		FatPer100:     off.Product.Nutriments.Fat100g,
		CarbsPer100:   off.Product.Nutriments.Carbohydrates100g,
		UnitType:      UnitWeight,
	}
	return p //возвращаем его
}
