package storage

type Storage interface {
	Districts() ([]District, error)
	PickDistrict(name string) (*District, error)
	PickPlace(name string) (*Place, error)
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
