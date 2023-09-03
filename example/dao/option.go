package dao

import (
	"gorm.io/gorm"
	"tccTrx/txManager"
)

type QueryOption func(db *gorm.DB) *gorm.DB

func WithID(id uint) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

func WithStatus(status txManager.ComponentTryStatus) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("status = ?", status.String())
	}
}
