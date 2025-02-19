package unitTests

import (
	"RemitlyTask/src/models"

	"github.com/stretchr/testify/mock"
)

type MockSwiftCodeRepository struct {
	mock.Mock
}

func (m *MockSwiftCodeRepository) FindBySwiftCodePrefix(prefix string) ([]models.SwiftCode, error) {
	args := m.Called(prefix)
	return args.Get(0).([]models.SwiftCode), args.Error(1)
}

func (m *MockSwiftCodeRepository) FindBySwiftCode(code string) (models.SwiftCode, error) {
	args := m.Called(code)
	return args.Get(0).(models.SwiftCode), args.Error(1)
}

func (m *MockSwiftCodeRepository) FindCountryNameByISO2(iso2 string) (string, error) {
	args := m.Called(iso2)
	return args.String(0), args.Error(1)
}

func (m *MockSwiftCodeRepository) FindByCountryISO2(iso2 string) ([]models.SwiftCode, error) {
	args := m.Called(iso2)
	return args.Get(0).([]models.SwiftCode), args.Error(1)
}

func (m *MockSwiftCodeRepository) Create(newCode *models.SwiftCode) error {
	args := m.Called(newCode)
	return args.Error(0)
}

func (m *MockSwiftCodeRepository) Delete(swiftCode string) error {
	args := m.Called(swiftCode)
	return args.Error(0)
}
