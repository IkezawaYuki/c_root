package filter

import "gorm.io/gorm"

type Filter interface {
	GenerateMods(db *gorm.DB) *gorm.DB
}
