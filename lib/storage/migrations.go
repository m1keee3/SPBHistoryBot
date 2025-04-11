package storage

import (
	"SPBHistoryBot/lib/e"
	"log"
)

func Init(db *DBStorage) error {
	if err := db.db.AutoMigrate(&District{}, &Place{}); err != nil {
		return e.Wrap("Failed to migrate models", err)
	}
	return nil
}

func SeedifEmpty(db *DBStorage) error {
	tables := []interface{}{&District{}, &Place{}}
	empty := true

	for _, t := range tables {
		var count int64
		if err := db.db.Model(t).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			empty = false
			break
		}
	}

	if empty {
		return Seed(db)
	}

	log.Print("skip seeding, because DB is not empty")
	return nil
}

func Seed(db *DBStorage) error {

	place := Place{
		Name: "Исаакиевский собор",
		Text: "Построен в 1818–1858 годах",
	}

	district := District{
		Name:   "Невский район",
		Text:   "Центральный район города",
		Places: []Place{place},
	}

	if err := db.db.Create(&district).Error; err != nil {
		return e.Wrap("Failed to seed data", err)
	}

	return nil
}
