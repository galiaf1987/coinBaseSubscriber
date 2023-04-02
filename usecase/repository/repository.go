package repository

import (
	"github.com/jinzhu/gorm"
)

type BaseRepository struct {
	DBConnection *gorm.DB
}

func (r BaseRepository) DB() *gorm.DB {
	return r.DBConnection
}
