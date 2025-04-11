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
	if err := db.db.Preload("Place").Find(&districts).Error; err != nil {
		return nil, e.Wrap("Failed to get districts", err)
	}
	return districts, nil
}

func (s *DBStorage) PickDistrict(name string) (*District, error) {
	var district District
	err := s.db.Preload("Places").Where("name = ?", name).First(&district).Error
	if err != nil {
		return nil, e.Wrap("Failed to get district", err)
	}
	return &district, nil
}

func (s *DBStorage) PickPlace(name string) (*Place, error) {
	var place Place
	err := s.db.Where("name = ?", name).First(&place).Error
	if err != nil {
		return nil, e.Wrap("Failed to get district place", err)
	}
	return &place, nil
}
