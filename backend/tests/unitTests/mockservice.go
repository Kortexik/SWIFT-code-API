package unitTests

import (
	"RemitlyTask/src/models"

	"github.com/stretchr/testify/mock"
)

type MockSwiftCodeService struct {
	mock.Mock
}

func (m *MockSwiftCodeService) AddSwiftCode(newCode *models.SwiftCode) error {
	args := m.Called(newCode)
	return args.Error(0)
}

func (m *MockSwiftCodeService) DeleteSwiftCode(swiftCode string) error {
	args := m.Called(swiftCode)
	return args.Error(0)
}

func (m *MockSwiftCodeService) GetSwiftCodesByCountry(iso2 string) (interface{}, error) {
	args := m.Called(iso2)
	return args.Get(0), args.Error(1)
}

func (m *MockSwiftCodeService) GetHeadquarterDetails(swiftCodePrefix string) (interface{}, error) {
	args := m.Called(swiftCodePrefix)
	return args.Get(0), args.Error(1)
}

func (m *MockSwiftCodeService) GetBranchDetails(swiftCode string) (interface{}, error) {
	args := m.Called(swiftCode)
	return args.Get(0), args.Error(1)
}

func (m *MockSwiftCodeService) GetCountryName(iso2 string) (string, error) {
	args := m.Called(iso2)
	return args.String(0), args.Error(1)
}
