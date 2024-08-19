package repository

import "gorm.io/gorm"

type BaseRepository struct {
	db *gorm.DB
}

func NewBaseRepository(db *gorm.DB) *BaseRepository {
	return &BaseRepository{db: db}
}

func (b *BaseRepository) Begin() *gorm.DB {
	return b.db.Begin()
}

func (b *BaseRepository) Commit(tx *gorm.DB) error {
	return tx.Commit().Error
}

func (b *BaseRepository) Rollback(tx *gorm.DB) error {
	return tx.Rollback().Error
}
