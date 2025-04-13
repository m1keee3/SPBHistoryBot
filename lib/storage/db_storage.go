package storage

import (
	"SPBHistoryBot/lib/e"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBStorage struct {
	db *gorm.DB
}

func NewDBStorage(dsn string) (*DBStorage, error) {

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, e.Wrap("Failed to connect to DB", err)
	}

	return &DBStorage{db: db}, nil
}

func (db *DBStorage) Districts() ([]District, error) {
	var districts []District
	if err := db.db.Preload("Places").Find(&districts).Error; err != nil {
		return nil, e.Wrap("Failed to get districts", err)
	}
	return districts, nil
}

func (s *DBStorage) FindDistrict(id int) (*District, error) {
	var district District
	err := s.db.Preload("Places").Where("id = ?", id).First(&district).Error
	if err != nil {
		return nil, e.Wrap("Failed to get district", err)
	}
	return &district, nil
}

func (s *DBStorage) FindPlace(id int) (*Place, error) {
	var place Place
	err := s.db.Where("id = ?", id).First(&place).Error
	if err != nil {
		return nil, e.Wrap("Failed to get district place", err)
	}
	return &place, nil
}
