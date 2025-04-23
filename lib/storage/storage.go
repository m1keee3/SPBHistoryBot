package storage

type Storage interface {
	Districts() ([]District, error)
	FindDistrict(id int) (*District, error)
	FindPlace(id int) (*Place, error)
	FindNearPlace(latitude float64, longitude float64) (*Place, error)
}

type Place struct {
	ID         uint `gorm:"primaryKey"`
	Name       string
	Text       string
	Image      string
	Latitude   float64
	Longitude  float64
	DistrictID uint
}

type District struct {
	ID     uint `gorm:"primaryKey"`
	Name   string
	Text   string
	Image  string
	Places []Place
}
