package model

import (
	"gorm.io/gorm"
)

type PV struct {
	gorm.Model

	PageKey  string `gorm:"index;size:255"` // 浏览页面
	SiteName string `gorm:"index;size:255"` // 网站

	Num uint // 浏览量
}

func (v *PV) IsEmpty() bool {
	return v.ID == 0
}
