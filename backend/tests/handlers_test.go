package tests

import (
	"RemitlyTask/src/handlers"
	"RemitlyTask/src/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAddSwiftCodeIntegration(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	handler := handlers.NewSwiftCodeHandler(db)

	r := gin.Default()
	vCodes := r.Group("v1/swift-codes")
	{
		vCodes.POST("/", handler.AddNewSwiftCode)
	}
	ts := httptest.NewServer(r)
	defer ts.Close()

	testCase := models.SwiftCodeBranch{
		Address:       "TEST ADDRESS",
		BankName:      "TEST BANK",
		CountryISO2:   "PL",
		CountryName:   "POLAND",
		IsHeadquarter: false,
		SwiftCode:     "TESTTESTTES",
	}
	expectedStatus := http.StatusOK
	expectedResponse := map[string]string{"message": testCase.SwiftCode + " has been added to database."}

	t.Run("Add new Swift code", func(t *testing.T) {
		payload, err := json.Marshal(testCase)
		assert.NoError(t, err)

		req, err := http.NewRequest("POST", ts.URL+"/v1/swift-codes/", bytes.NewBuffer(payload))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, expectedStatus, resp.StatusCode)

		var response map[string]string
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
	})
}

func TestGetCodeIntegration(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	handler := handlers.NewSwiftCodeHandler(db)

	r := gin.Default()
	vCodes := r.Group("v1/swift-codes")
	{
		vCodes.GET("/:swift-code", handler.GetCode)
	}
	ts := httptest.NewServer(r)
	defer ts.Close()

	testCase := models.SwiftCode{
		Address:     "TEST ADDRESS",
		Name:        "TEST BANK",
		CountryISO2: "PL",
		CountryName: "POLAND",
		SwiftCode:   "TESTTESTTES",
	}

	testCase2 := models.SwiftCode{
		Address:     "TEST ADDRESS HQ",
		Name:        "TEST HQ BANK",
		CountryISO2: "PL",
		CountryName: "POLAND",
		SwiftCode:   "TESTTESTXXX",
	}

	err := db.Create(&testCase).Error
	assert.NoError(t, err)
	err = db.Create(&testCase2).Error
	assert.NoError(t, err)

	expectedStatus := http.StatusOK
	expectedResponse := models.SwiftCodeBranch{
		Address:       testCase.Address,
		BankName:      testCase.Name,
		CountryISO2:   testCase.CountryISO2,
		CountryName:   testCase.CountryName,
		IsHeadquarter: false,
		SwiftCode:     testCase.SwiftCode,
	}
	expectedResponse2 := models.SwiftCodeDetails{
		Address:       testCase2.Address,
		BankName:      testCase2.Name,
		CountryISO2:   testCase2.CountryISO2,
		CountryName:   testCase2.CountryName,
		IsHeadquarter: true,
		SwiftCode:     testCase2.SwiftCode,
		Branches: []models.SwiftCodeBank{
			{
				Address:       testCase.Address,
				BankName:      testCase.Name,
				CountryISO2:   testCase.CountryISO2,
				IsHeadquarter: false,
				SwiftCode:     testCase.SwiftCode,
			},
		},
	}

	t.Run("Get details by swift-code", func(t *testing.T) {
		req, err := http.NewRequest("GET", ts.URL+"/v1/swift-codes/"+testCase.SwiftCode, nil)
		assert.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, expectedStatus, resp.StatusCode)

		var response models.SwiftCodeBranch
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
		//---------------------------------------------------
		req2, err := http.NewRequest("GET", ts.URL+"/v1/swift-codes/"+testCase2.SwiftCode, nil)
		assert.NoError(t, err)

		resp2, err := http.DefaultClient.Do(req2)
		assert.NoError(t, err)
		defer resp2.Body.Close()

		assert.Equal(t, expectedStatus, resp2.StatusCode)

		var response2 models.SwiftCodeDetails
		err = json.NewDecoder(resp2.Body).Decode(&response2)
		assert.NoError(t, err)
		assert.Equal(t, expectedResponse2, response2)
	})
}

func TestGetCodesByCountry(t *testing.T) {

	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	handler := handlers.NewSwiftCodeHandler(db)

	r := gin.Default()
	vCodes := r.Group("v1/swift-codes")
	{
		vCodes.GET("/country/:ISO2", handler.GetCodesByCountry)
	}
	ts := httptest.NewServer(r)
	defer ts.Close()

	testCase := models.SwiftCode{
		Address:     "TEST ADDRESS",
		Name:        "TEST BANK",
		CountryISO2: "PL",
		CountryName: "POLAND",
		SwiftCode:   "TESTTESTTES",
	}

	testCase2 := models.SwiftCode{
		Address:     "TEST ADDRESS HQ",
		Name:        "TEST HQ BANK",
		CountryISO2: "PL",
		CountryName: "POLAND",
		SwiftCode:   "TESTTESTXXX",
	}

	err := db.Create(&testCase).Error
	assert.NoError(t, err)
	err = db.Create(&testCase2).Error
	assert.NoError(t, err)

	var banks []models.SwiftCodeBank
	bank1 := models.SwiftCodeBank{
		Address:       testCase.Address,
		BankName:      testCase.Name,
		CountryISO2:   testCase.CountryISO2,
		IsHeadquarter: false,
		SwiftCode:     testCase.SwiftCode,
	}
	bank2 := models.SwiftCodeBank{
		Address:       testCase2.Address,
		BankName:      testCase2.Name,
		CountryISO2:   testCase2.CountryISO2,
		IsHeadquarter: true,
		SwiftCode:     testCase2.SwiftCode,
	}

	banks = append(banks, bank1, bank2)

	expectedStatus := http.StatusOK
	expectedResponse := models.SwiftCodeCountry{
		CountryISO2: testCase.CountryISO2,
		CountryName: testCase.CountryName,
		SwiftCodes:  banks,
	}
	var ISO2 = "PL"
	t.Run("Get swift codes by country", func(t *testing.T) {
		req, err := http.NewRequest("GET", ts.URL+"/v1/swift-codes/country/"+ISO2, nil)
		assert.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, expectedStatus, resp.StatusCode)

		var response models.SwiftCodeCountry
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
	})
}

func TestDeleteCode(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	handler := handlers.NewSwiftCodeHandler(db)

	r := gin.Default()
	vCodes := r.Group("v1/swift-codes")
	{
		vCodes.DELETE("/:swift-code", handler.DeleteCode)
	}
	ts := httptest.NewServer(r)
	defer ts.Close()

	testCase := models.SwiftCode{
		Address:     "TEST ADDRESS",
		Name:        "TEST BANK",
		CountryISO2: "PL",
		CountryName: "POLAND",
		SwiftCode:   "TESTTESTTES",
	}

	err := db.Create(&testCase).Error
	assert.NoError(t, err)

	expectedStatus := http.StatusOK
	expectedResponse := map[string]string{"message": testCase.SwiftCode + " was removed."}

	t.Run("Delete swift code", func(t *testing.T) {
		req, err := http.NewRequest("DELETE", ts.URL+"/v1/swift-codes/"+testCase.SwiftCode, nil)
		assert.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, expectedStatus, resp.StatusCode)

		var response map[string]string
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
	})

}
