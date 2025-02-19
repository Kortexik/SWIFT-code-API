package unitTests

import (
	"RemitlyTask/src/handlers"
	"RemitlyTask/src/models"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAddNewSwiftCode(t *testing.T) {
	mockService := new(MockSwiftCodeService)
	handler := handlers.NewSwiftCodeHandlerByService(mockService)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/swift-codes", handler.AddNewSwiftCode)

	validCode := &models.SwiftCodeBranch{
		Address:       "123 Test St",
		BankName:      "Test Bank HQ",
		CountryISO2:   "PL",
		CountryName:   "POLAND",
		IsHeadquarter: false,
		SwiftCode:     "TESTUSABXYZ",
	}

	t.Run("TestAddNewSwiftCode_successful", func(t *testing.T) {
		mockService.On("GetCountryName", "PL").Return("POLAND", nil)
		mockService.On("GetBranchDetails", "TESTUSABXYZ").Return(nil, nil)
		mockService.On("AddSwiftCode", mock.AnythingOfType("*models.SwiftCode")).Return(nil)

		jsonData, err := json.Marshal(validCode)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/swift-codes", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]string
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, validCode.SwiftCode+" has been added to the database.", response["message"])
		mockService.AssertExpectations(t)
	})

	t.Run("TestAddNewSwiftCode_invalidIso2Length", func(t *testing.T) {
		invalidCode := &models.SwiftCodeBranch{
			Address:       "123 Test St",
			BankName:      "Test Bank HQ",
			CountryISO2:   "PLN",
			CountryName:   "POLAND",
			IsHeadquarter: false,
			SwiftCode:     "TESTUSABXYZ",
		}

		jsonData, err := json.Marshal(invalidCode)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/swift-codes", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("TestAddNewSwiftCode_invalidJSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/swift-codes", bytes.NewBuffer([]byte(`{"invalid json"}`)))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGetCode(t *testing.T) {
	mockService := new(MockSwiftCodeService)
	handler := handlers.NewSwiftCodeHandlerByService(mockService)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/swift-codes/:swift-code", handler.GetCode)

	branchCode := &models.SwiftCode{
		SwiftCode:   "TESTUSABXYZ",
		Name:        "Test Bank Branch",
		CountryISO2: "PL",
		Address:     "123 Test St",
		CountryName: "POLAND",
	}

	headquarterCode := &models.SwiftCode{
		SwiftCode:   "TESTUSABXXX",
		Name:        "Test Bank HQ",
		CountryISO2: "PL",
		Address:     "456 Test St",
		CountryName: "POLAND",
	}

	branches := []models.SwiftCodeBank{
		{
			SwiftCode:     branchCode.SwiftCode,
			BankName:      branchCode.Name,
			CountryISO2:   branchCode.CountryISO2,
			Address:       branchCode.Address,
			IsHeadquarter: false,
		},
	}

	headquarterResponse := models.SwiftCodeDetails{
		Address:       headquarterCode.Address,
		BankName:      headquarterCode.Name,
		CountryISO2:   headquarterCode.CountryISO2,
		CountryName:   headquarterCode.CountryName,
		IsHeadquarter: true,
		SwiftCode:     headquarterCode.SwiftCode,
		Branches:      branches,
	}

	t.Run("TestGetCode_invalidSwiftCodeFormat", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/swift-codes/INVALID", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("TestGetCode_branchSuccessful", func(t *testing.T) {
		mockService.On("GetBranchDetails", branchCode.SwiftCode).Return(branchCode, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/swift-codes/"+branchCode.SwiftCode, nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.SwiftCode
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, *branchCode, response)
		mockService.AssertExpectations(t)
	})

	t.Run("TestGetCode_headquartersSuccessful", func(t *testing.T) {
		mockService.On("GetHeadquarterDetails", headquarterCode.SwiftCode[:8]).Return(headquarterResponse, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/swift-codes/"+headquarterCode.SwiftCode, nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.SwiftCodeDetails
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, headquarterResponse, response)
		mockService.AssertExpectations(t)
	})

	t.Run("TestGetCode_notFound", func(t *testing.T) {
		notFoundCode := "NOTFOUNDXXX"

		mockService.On("GetHeadquarterDetails", notFoundCode[:8]).Return(nil, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/swift-codes/"+notFoundCode, nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "No SWIFT code found for: "+notFoundCode, response["message"])
		mockService.AssertExpectations(t)
	})
}

func TestGetCodesByCountry(t *testing.T) {
	mockService := new(MockSwiftCodeService)
	handler := handlers.NewSwiftCodeHandlerByService(mockService)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/swift-codes/country/:ISO2", handler.GetCodesByCountry)

	expectedResponse := models.SwiftCodeCountry{
		CountryISO2: "PL",
		CountryName: "POLAND",
		SwiftCodes: []models.SwiftCodeBank{
			{
				SwiftCode:     "TESTUSABXXX",
				BankName:      "Test Bank HQ",
				CountryISO2:   "PL",
				Address:       "123 Test St",
				IsHeadquarter: true,
			},
			{
				SwiftCode:     "TESTUSABXYZ",
				BankName:      "Test Bank Branch",
				CountryISO2:   "PL",
				Address:       "456 Test St",
				IsHeadquarter: false,
			},
		},
	}

	t.Run("TestGetCodesByCountry_successful", func(t *testing.T) {
		mockService.On("GetSwiftCodesByCountry", "PL").Return(expectedResponse, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/swift-codes/country/PL", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.SwiftCodeCountry
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, response)

		mockService.AssertExpectations(t)
	})

	t.Run("TestGetCodesByCountry_invalidISO2Length", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/swift-codes/country/USA", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["message"], "Invalid ISO2 code")
	})

	t.Run("TestGetCodesByCountry_NonExistentISO2", func(t *testing.T) {
		emptyResponse := models.SwiftCodeCountry{
			CountryISO2: "XX",
			CountryName: "",
			SwiftCodes:  []models.SwiftCodeBank{},
		}
		mockService.On("GetSwiftCodesByCountry", "XX").Return(emptyResponse, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/swift-codes/country/XX", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("TestGetCodesByCountry_serviceError", func(t *testing.T) {
		mockService.On("GetSwiftCodesByCountry", "FR").Return(models.SwiftCodeCountry{}, errors.New("service error"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/swift-codes/country/FR", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestDeleteCode(t *testing.T) {
	mockService := new(MockSwiftCodeService)
	handler := handlers.NewSwiftCodeHandlerByService(mockService)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.DELETE("/swift-codes/:swift-code", handler.DeleteCode)

	t.Run("TestDeleteCode_successful", func(t *testing.T) {
		swiftCode := "TESTTESTXXX"
		mockService.On("DeleteSwiftCode", swiftCode).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/swift-codes/"+swiftCode, nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, swiftCode+" was removed.", response["message"])
		mockService.AssertExpectations(t)
	})

	t.Run("TestDeleteCode_nonExistentCode", func(t *testing.T) {
		swiftCode := "NONEXISTXXX"
		mockService.On("DeleteSwiftCode", swiftCode).Return(errors.New("SWIFT code NONEXISTENT not found"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/swift-codes/"+swiftCode, nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Could not delete a record SWIFT code NONEXISTENT not found", response["message"])
		mockService.AssertExpectations(t)
	})
}
