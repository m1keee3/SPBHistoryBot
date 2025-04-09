package storage

type Storage interface {
	Save()
	PickDistrict(name string) (*District, error)
	PickPlace(name string) (*Place, error)
}

type Place struct {
	Name string
	Text string
}

type District struct {
	Name   string
	Text   string
	Places []Place
}
