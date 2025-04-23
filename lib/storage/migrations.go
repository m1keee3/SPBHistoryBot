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

	peter_fort := Place{
		Name:      "Петропавловская крепость",
		Text:      "Форма звезды",
		Image:     "https://static.78.ru/images/uploads/1686542964143.jpg",
		Latitude:  59.9500019,
		Longitude: 30.3166718,
	}

	is_sobor := Place{
		Name:      "Исаакиевский собор",
		Text:      "Построен в 1818–1858 годах",
		Image:     "https://yandex.ru/images/search?pos=5&from=tabbar&img_url=https%3A%2F%2Fsneg.top%2Fuploads%2Fposts%2F2023-03%2F1680058387_sneg-top-p-isaakievskii-sobor-oboi-na-telefon-pintere-26.jpg&text=санкт-петербург&rpt=simage&lr=2",
		Latitude:  59.934111,
		Longitude: 30.306074,
	}
	kaz_sobor := Place{
		Name:      "Казанский собор",
		Text:      "Полукруглый",
		Image:     "https://ic.pics.livejournal.com/air_vision/76938106/13174/13174_original.jpg",
		Latitude:  59.93424,
		Longitude: 30.32444,
	}
	spas_na_kr := Place{
		Name:      "Спас на Крови",
		Text:      "Красивый",
		Image:     "https://image.fonwall.ru/o/6z/hram-spas-na-krovi-sankt-peterburg-rossiya.jpg?auto=compress&amp;fit=crop&amp;w=2560&amp;h=1440&amp;domain=img.fonwall.ru",
		Latitude:  59.940055,
		Longitude: 30.328807,
	}
	gen_shtab := Place{
		Name:      "Главный штаб",
		Text:      "Большой",
		Image:     "https://wallpapers.com/images/hd/russia-general-staff-building-ub6k2x9y4ifuadfp.jpg",
		Latitude:  59.933973,
		Longitude: 30.306103,
	}

	nevskiy := District{
		Name:   "Невский район",
		Text:   "Центральный район города",
		Places: []Place{peter_fort, is_sobor, kaz_sobor, spas_na_kr, gen_shtab},
	}

	if err := db.db.Create(&nevskiy).Error; err != nil {
		return e.Wrap("Failed to seed data", err)
	}

	return nil
}
