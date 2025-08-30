package model

import "time"

type Value struct {
	String_value  *string
	Int64_value   *int64
	Float64_value *float64
	Bool_value    *bool
}

type Category int

const (
	CategoryUnknown  Category = iota // 0 — Неизвестная категория
	CategoryEngine                   // 1 — Двигатель
	CategoryFuel                     // 2 — Топливо
	CategoryPorthole                 // 3 — Иллюминатор
	CategoryWing                     // 4 — Крыло
)

type Dimensions struct {
	Length float64
	Width  float64
	Height float64
	Weight float64
}

// NewDimensions creates a new Dimensions struct
func NewDimensions(length, width, height, weight float64) *Dimensions {
	return &Dimensions{
		Length: length,
		Width:  width,
		Height: height,
		Weight: weight,
	}
}

type Manufacturer struct {
	Name    string
	Country string
	Website string
}

// NewManufacturer creates a new Manufacturer struct
func NewManufacturer(name, country, website string) *Manufacturer {
	return &Manufacturer{
		Name:    name,
		Country: country,
		Website: website,
	}
}

type PartInfo struct {
	Name          string
	Description   string
	Price         float64
	StockQuantity int64
	Tags          []string
	Category      Category
	Dimensions    *Dimensions
	Manufacturer  *Manufacturer
	Metadata      map[string]*Value
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
}

type Part struct {
	Uuid string
	Info PartInfo
}
