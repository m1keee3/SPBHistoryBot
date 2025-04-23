package storage

import (
	"SPBHistoryBot/lib/e"
	"errors"
	"github.com/umahmood/haversine"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"math"
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

func (db *DBStorage) FindDistrict(id int) (*District, error) {
	var district District
	err := db.db.Preload("Places").Where("id = ?", id).First(&district).Error
	if err != nil {
		return nil, e.Wrap("Failed to get district", err)
	}
	return &district, nil
}

func (db *DBStorage) FindPlace(id int) (*Place, error) {
	var place Place
	err := db.db.Where("id = ?", id).First(&place).Error
	if err != nil {
		return nil, e.Wrap("Failed to get district place", err)
	}
	return &place, nil
}

func (db *DBStorage) FindNearPlace(latitude float64, longitude float64) (*Place, error) {
	var places []Place
	if err := db.db.Find(&places).Error; err != nil {
		return nil, e.Wrap("failed to load places", err)
	}

	if len(places) == 0 {
		return nil, errors.New("no places found")
	}

	origin := haversine.Coord{Lat: latitude, Lon: longitude}

	var nearest *Place
	minDistance := math.MaxFloat64

	for i, p := range places {
		destination := haversine.Coord{Lat: p.Latitude, Lon: p.Longitude}
		km, _ := haversine.Distance(origin, destination)
		if km < minDistance {
			minDistance = km
			nearest = &places[i]
		}
	}

	return nearest, nil
}
