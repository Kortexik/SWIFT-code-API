package unitTests

import (
	"RemitlyTask/src/models"
	"RemitlyTask/src/services"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHeadquarterDetails(t *testing.T) {
	t.Run("TestGetHeadquarterDetails_successful", func(t *testing.T) {
		mockRepo := &MockSwiftCodeRepository{}
		service := services.NewSwiftCodeService(mockRepo)

		swiftCodes := []models.SwiftCode{
			{
				SwiftCode:   "TESTUSABXXX",
				Name:        "Test Bank HQ",
				CountryISO2: "US",
				Address:     "123 Test St",
				CountryName: "UNITED STATES",
			},
			{
				SwiftCode:   "TESTUSABNYC",
				Name:        "Test Bank Branch",
				CountryISO2: "US",
				Address:     "456 Branch St",
				CountryName: "UNITED STATES",
			},
		}
		mockRepo.On("FindBySwiftCodePrefix", "TESTUSAB").Return(swiftCodes, nil)

		response, err := service.GetHeadquarterDetails("TESTUSAB")

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, "Test Bank HQ", response.(models.SwiftCodeDetails).BankName)
		assert.Equal(t, "123 Test St", response.(models.SwiftCodeDetails).Address)
		assert.Equal(t, "US", response.(models.SwiftCodeDetails).CountryISO2)
		assert.Equal(t, "UNITED STATES", response.(models.SwiftCodeDetails).CountryName)
		assert.True(t, response.(models.SwiftCodeDetails).IsHeadquarter)
		assert.Len(t, response.(models.SwiftCodeDetails).Branches, 1)
		mockRepo.AssertExpectations(t)
	})

	t.Run("TestGetHeadquarterDetails_noHeadquarterFound", func(t *testing.T) {
		mockRepo := &MockSwiftCodeRepository{}
		service := services.NewSwiftCodeService(mockRepo)

		swiftCodes := []models.SwiftCode{
			{
				SwiftCode:   "TESTUSABNYC",
				Name:        "Test Bank Branch",
				CountryISO2: "US",
				Address:     "456 Branch St",
				CountryName: "UNITED STATES",
			},
		}
		mockRepo.On("FindBySwiftCodePrefix", "TESTUSAB").Return(swiftCodes, nil)

		response, err := service.GetHeadquarterDetails("TESTUSAB")

		assert.NoError(t, err)
		assert.Nil(t, response)
		mockRepo.AssertExpectations(t)
	})

	t.Run("TestGetHeadquarterDetails_repositoryError", func(t *testing.T) {
		mockRepo := &MockSwiftCodeRepository{}
		service := services.NewSwiftCodeService(mockRepo)

		mockRepo.On("FindBySwiftCodePrefix", "TESTUSAB").Return([]models.SwiftCode{}, errors.New("repository error"))

		response, err := service.GetHeadquarterDetails("TESTUSAB")

		assert.Error(t, err)
		assert.Nil(t, response)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetBranchDetails(t *testing.T) {
	t.Run("TestGetBranchDetails_successful", func(t *testing.T) {
		mockRepo := &MockSwiftCodeRepository{}
		service := services.NewSwiftCodeService(mockRepo)

		branch := models.SwiftCode{
			SwiftCode:   "TESTUSABNYC",
			Name:        "Test Bank Branch",
			CountryISO2: "US",
			Address:     "456 Branch St",
			CountryName: "UNITED STATES",
		}
		mockRepo.On("FindBySwiftCode", "TESTUSABNYC").Return(branch, nil)

		response, err := service.GetBranchDetails("TESTUSABNYC")

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, "Test Bank Branch", response.(models.SwiftCodeBranch).BankName)
		assert.Equal(t, "456 Branch St", response.(models.SwiftCodeBranch).Address)
		assert.Equal(t, "US", response.(models.SwiftCodeBranch).CountryISO2)
		assert.Equal(t, "UNITED STATES", response.(models.SwiftCodeBranch).CountryName)
		assert.False(t, response.(models.SwiftCodeBranch).IsHeadquarter)
		mockRepo.AssertExpectations(t)
	})

	t.Run("TestGetBranchDetails_branchNotFound", func(t *testing.T) {
		mockRepo := &MockSwiftCodeRepository{}
		service := services.NewSwiftCodeService(mockRepo)

		mockRepo.On("FindBySwiftCode", "TESTUSABNYC").Return(models.SwiftCode{}, nil)

		response, err := service.GetBranchDetails("TESTUSABNYC")

		assert.NoError(t, err)
		assert.Nil(t, response)
		mockRepo.AssertExpectations(t)
	})

	t.Run("TestGetBranchDetails_repositoryError", func(t *testing.T) {
		mockRepo := &MockSwiftCodeRepository{}
		service := services.NewSwiftCodeService(mockRepo)

		mockRepo.On("FindBySwiftCode", "TESTUSABNYC").Return(models.SwiftCode{}, errors.New("Repository error"))

		response, err := service.GetBranchDetails("TESTUSABNYC")

		assert.Error(t, err)
		assert.Nil(t, response)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetSwiftCodesByCountry(t *testing.T) {
	t.Run("TestGetSwiftCodesByCountry_successful", func(t *testing.T) {
		mockRepo := &MockSwiftCodeRepository{}
		service := services.NewSwiftCodeService(mockRepo)

		swiftCodes := []models.SwiftCode{
			{
				SwiftCode:   "TESTUSABXXX",
				Name:        "Test Bank HQ",
				CountryISO2: "US",
				Address:     "123 Test St",
				CountryName: "UNITED STATES",
			},
			{
				SwiftCode:   "TESTUSABNYC",
				Name:        "Test Bank Branch",
				CountryISO2: "US",
				Address:     "456 Branch St",
				CountryName: "UNITED STATES",
			},
		}
		mockRepo.On("FindCountryNameByISO2", "US").Return("UNITED STATES", nil)
		mockRepo.On("FindByCountryISO2", "US").Return(swiftCodes, nil)

		response, err := service.GetSwiftCodesByCountry("US")

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, "US", response.(models.SwiftCodeCountry).CountryISO2)
		assert.Equal(t, "UNITED STATES", response.(models.SwiftCodeCountry).CountryName)
		assert.Len(t, response.(models.SwiftCodeCountry).SwiftCodes, 2)
		mockRepo.AssertExpectations(t)
	})

	t.Run("TestGetSwiftCodesByCountry_noCodesFoundForCountry", func(t *testing.T) {
		mockRepo := &MockSwiftCodeRepository{}
		service := services.NewSwiftCodeService(mockRepo)

		mockRepo.On("FindCountryNameByISO2", "US").Return("UNITED STATES", nil)
		mockRepo.On("FindByCountryISO2", "US").Return([]models.SwiftCode{}, nil)

		response, err := service.GetSwiftCodesByCountry("US")

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, "US", response.(models.SwiftCodeCountry).CountryISO2)
		assert.Equal(t, "UNITED STATES", response.(models.SwiftCodeCountry).CountryName)
		assert.Len(t, response.(models.SwiftCodeCountry).SwiftCodes, 0)
		mockRepo.AssertExpectations(t)
	})

	t.Run("TestGetSwiftCodesByCountry_repositoryError", func(t *testing.T) {
		mockRepo := &MockSwiftCodeRepository{}
		service := services.NewSwiftCodeService(mockRepo)

		mockRepo.On("FindCountryNameByISO2", "US").Return("", errors.New("repository error"))

		response, err := service.GetSwiftCodesByCountry("US")

		assert.Error(t, err)
		assert.Nil(t, response)
		mockRepo.AssertExpectations(t)
	})
}

func TestAddSwiftCode(t *testing.T) {
	t.Run("TestAddSwiftCode_successful", func(t *testing.T) {
		mockRepo := &MockSwiftCodeRepository{}
		service := services.NewSwiftCodeService(mockRepo)

		newCode := &models.SwiftCode{
			SwiftCode:   "TESTUSABXXX",
			Name:        "Test Bank HQ",
			CountryISO2: "US",
			Address:     "123 Test St",
			CountryName: "UNITED STATES",
		}
		mockRepo.On("Create", newCode).Return(nil)

		err := service.AddSwiftCode(newCode)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("TestAddSwiftCode_repositoryError", func(t *testing.T) {
		mockRepo := &MockSwiftCodeRepository{}
		service := services.NewSwiftCodeService(mockRepo)

		newCode := &models.SwiftCode{
			SwiftCode:   "TESTUSABXXX",
			Name:        "Test Bank HQ",
			CountryISO2: "US",
			Address:     "123 Test St",
			CountryName: "UNITED STATES",
		}
		mockRepo.On("Create", newCode).Return(errors.New("repository error"))

		err := service.AddSwiftCode(newCode)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteSwiftCode(t *testing.T) {
	t.Run("TestDeleteSwiftCode_successful", func(t *testing.T) {
		mockRepo := &MockSwiftCodeRepository{}
		service := services.NewSwiftCodeService(mockRepo)

		mockRepo.On("Delete", "TESTUSABXXX").Return(nil)

		err := service.DeleteSwiftCode("TESTUSABXXX")

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("TestDeleteSwiftCode_repositoryError", func(t *testing.T) {
		mockRepo := &MockSwiftCodeRepository{}
		service := services.NewSwiftCodeService(mockRepo)

		mockRepo.On("Delete", "TESTUSABXXX").Return(errors.New("Repository error"))

		err := service.DeleteSwiftCode("TESTUSABXXX")

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
