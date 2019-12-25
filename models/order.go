package models

import "github.com/jinzhu/gorm"

const (
	Open int = iota
	Placed
	Received
)

type Order struct {
	gorm.Model
	Vendor      string
	PartNumber  string
	Description string
	Quantity    int
	UnitCost    float32
	Status      int
}
