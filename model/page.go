package model

import "gorm.io/gorm"

type Page struct {
	gorm.Model
	Name string
	Link string
}
