package storage

type Storage interface {
	Districts() ([]District, error)
	PickDistrict(id int) (*District, error)
	PickPlace(id int) (*Place, error)
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
