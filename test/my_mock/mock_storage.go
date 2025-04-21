package my_mock

import (
	"SPBHistoryBot/lib/storage"
	"github.com/stretchr/testify/mock"
)

type Storage struct {
	mock.Mock
}

func (m *Storage) Districts() ([]storage.District, error) {
	args := m.Called()
	return args.Get(0).([]storage.District), args.Error(1)
}

func (m *Storage) FindDistrict(id int) (*storage.District, error) {
	args := m.Called(id)
	return args.Get(0).(*storage.District), args.Error(1)
}

func (m *Storage) FindPlace(id int) (*storage.Place, error) {
	args := m.Called(id)
	return args.Get(0).(*storage.Place), args.Error(1)
}
