package integrationTests

import (
	"RemitlyTask/src/handlers"
	"RemitlyTask/src/models"
	"RemitlyTask/tests/testHelpers"
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddSwiftCodeIntegration(t *testing.T) {
	db := testHelpers.SetupTestDB(t)
	defer testHelpers.CleanupTestDB(t, db)

	handler := handlers.NewSwiftCodeHandler(db)
	require.NotNil(t, handler)

	ts := testHelpers.SetupTestServer(handler)
	require.NotNil(t, ts)
	defer ts.Close()

	testCases := []struct {
		name             string
		input            models.SwiftCodeBranch
		expectedStatus   int
		expectedResponse map[string]string
	}{
		{
			name: "Valid SWIFT code",
			input: models.SwiftCodeBranch{
				Address:       "TEST ADDRESS",
				BankName:      "TEST BANK",
				CountryISO2:   "PL",
				CountryName:   "POLAND",
				IsHeadquarter: false,
				SwiftCode:     "TESTTESTTES",
			},
			expectedStatus: http.StatusOK,
			expectedResponse: map[string]string{
				"message": "TESTTESTTES has been added to the database.",
			},
		},
		{
			name: "Invalid SWIFT code length",
			input: models.SwiftCodeBranch{
				Address:       "TEST ADDRESS",
				BankName:      "TEST BANK",
				CountryISO2:   "PL",
				CountryName:   "POLAND",
				IsHeadquarter: false,
				SwiftCode:     "TESTTEE",
			},
			expectedStatus: http.StatusBadRequest,
			expectedResponse: map[string]string{
				"message": "Error inserting to database TESTTEE swift code must be 8 or 11 characters long.",
			},
		},
		{
			name: "Invalid ISO2 code length",
			input: models.SwiftCodeBranch{
				Address:       "TEST ADDRESS",
				BankName:      "TEST BANK",
				CountryISO2:   "P",
				CountryName:   "POLAND",
				IsHeadquarter: false,
				SwiftCode:     "TESTTESTTES",
			},
			expectedStatus: http.StatusBadRequest,
			expectedResponse: map[string]string{
				"message": "Error inserting to database this iso2 code: P does not exist.",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			payload, err := json.Marshal(tc.input)
			assert.NoError(t, err)

			req, err := http.NewRequest("POST", ts.URL+"/v1/swift-codes/", bytes.NewBuffer(payload))
			assert.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			var response map[string]string
			err = json.NewDecoder(resp.Body).Decode(&response)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, response)
		})
	}

}

func TestGetCodeIntegration(t *testing.T) {
	db := testHelpers.SetupTestDB(t)
	defer testHelpers.CleanupTestDB(t, db)

	handler := handlers.NewSwiftCodeHandler(db)
	require.NotNil(t, handler)

	ts := testHelpers.SetupTestServer(handler)
	require.NotNil(t, ts)
	defer ts.Close()

	testData := []models.SwiftCode{
		{
			Address:     "TEST ADDRESS",
			Name:        "TEST BANK",
			CountryISO2: "PL",
			CountryName: "POLAND",
			SwiftCode:   "TESTTESTTES",
		},
		{
			Address:     "TEST ADDRESS HQ",
			Name:        "TEST HQ BANK",
			CountryISO2: "PL",
			CountryName: "POLAND",
			SwiftCode:   "TESTTESTXXX",
		},
	}

	for _, data := range testData {
		err := db.Create(&data).Error
		assert.NoError(t, err)
	}

	testCases := []struct {
		name             string
		input            string
		expectedStatus   int
		expectedResponse interface{}
	}{
		{
			name:           "Get details by SWIFT code (Branch)",
			input:          "TESTTESTTES",
			expectedStatus: http.StatusOK,
			expectedResponse: models.SwiftCodeBranch{
				Address:       "TEST ADDRESS",
				BankName:      "TEST BANK",
				CountryISO2:   "PL",
				CountryName:   "POLAND",
				IsHeadquarter: false,
				SwiftCode:     "TESTTESTTES",
			},
		},
		{
			name:           "Get details by SWIFT code (HQ with branches)",
			input:          "TESTTESTXXX",
			expectedStatus: http.StatusOK,
			expectedResponse: models.SwiftCodeDetails{
				Address:       "TEST ADDRESS HQ",
				BankName:      "TEST HQ BANK",
				CountryISO2:   "PL",
				CountryName:   "POLAND",
				IsHeadquarter: true,
				SwiftCode:     "TESTTESTXXX",
				Branches: []models.SwiftCodeBank{
					{
						Address:       "TEST ADDRESS",
						BankName:      "TEST BANK",
						CountryISO2:   "PL",
						IsHeadquarter: false,
						SwiftCode:     "TESTTESTTES",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", ts.URL+"/v1/swift-codes/"+tc.input, nil)
			assert.NoError(t, err)

			resp, err := http.DefaultClient.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			if strings.HasSuffix(tc.input, "XXX") {
				var response models.SwiftCodeDetails
				err = json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResponse, response)
			} else {
				var response models.SwiftCodeBranch
				err = json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResponse, response)
			}

		})
	}
}

func TestGetCodesByCountry(t *testing.T) {
	db := testHelpers.SetupTestDB(t)
	defer testHelpers.CleanupTestDB(t, db)

	handler := handlers.NewSwiftCodeHandler(db)
	require.NotNil(t, handler)

	ts := testHelpers.SetupTestServer(handler)
	require.NotNil(t, ts)
	defer ts.Close()

	testData := []models.SwiftCode{
		{
			Address:     "TEST ADDRESS",
			Name:        "TEST BANK",
			CountryISO2: "PL",
			CountryName: "POLAND",
			SwiftCode:   "TESTTESTTES",
		},
		{
			Address:     "TEST ADDRESS HQ",
			Name:        "TEST HQ BANK",
			CountryISO2: "PL",
			CountryName: "POLAND",
			SwiftCode:   "TESTTESTXXX",
		},
	}

	for _, data := range testData {
		err := db.Create(&data).Error
		assert.NoError(t, err)
	}

	testCases := []struct {
		name             string
		input            string
		expectedStatus   int
		expectedResponse interface{}
	}{
		{
			name:           "Get SWIFT codes by valid ISO2 code",
			input:          "PL",
			expectedStatus: http.StatusOK,
			expectedResponse: models.SwiftCodeCountry{
				CountryISO2: "PL",
				CountryName: "POLAND",
				SwiftCodes: []models.SwiftCodeBank{
					{
						Address:       "TEST ADDRESS",
						BankName:      "TEST BANK",
						CountryISO2:   "PL",
						IsHeadquarter: false,
						SwiftCode:     "TESTTESTTES",
					},
					{
						Address:       "TEST ADDRESS HQ",
						BankName:      "TEST HQ BANK",
						CountryISO2:   "PL",
						IsHeadquarter: true,
						SwiftCode:     "TESTTESTXXX",
					},
				},
			},
		},
		{
			name:           "Get SWIFT codes by invalid ISO2 code",
			input:          "INVALID",
			expectedStatus: http.StatusBadRequest,
			expectedResponse: map[string]string{
				"message": "Invalid ISO2 code length. It must be 2 characters long.",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", ts.URL+"/v1/swift-codes/country/"+tc.input, nil)
			assert.NoError(t, err)

			resp, err := http.DefaultClient.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			if tc.expectedStatus == http.StatusOK {
				var response models.SwiftCodeCountry
				err = json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResponse, response)
			} else {
				var response map[string]string
				err = json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResponse, response)
			}
		})
	}
}

func TestDeleteCode(t *testing.T) {
	db := testHelpers.SetupTestDB(t)
	defer testHelpers.CleanupTestDB(t, db)

	handler := handlers.NewSwiftCodeHandler(db)
	require.NotNil(t, handler)

	ts := testHelpers.SetupTestServer(handler)
	require.NotNil(t, ts)
	defer ts.Close()

	testData := models.SwiftCode{
		Address:     "TEST ADDRESS",
		Name:        "TEST BANK",
		CountryISO2: "PL",
		CountryName: "POLAND",
		SwiftCode:   "TESTTESTTES",
	}

	err := db.Create(&testData).Error
	assert.NoError(t, err)

	testCases := []struct {
		name             string
		input            string
		expectedStatus   int
		expectedResponse map[string]string
	}{
		{
			name:           "Delete existing SWIFT code",
			input:          "TESTTESTTES",
			expectedStatus: http.StatusOK,
			expectedResponse: map[string]string{
				"message": "TESTTESTTES was removed.",
			},
		},
		{
			name:           "Delete non-existent SWIFT code",
			input:          "NONEXISTENT",
			expectedStatus: http.StatusNotFound,
			expectedResponse: map[string]string{
				"message": "Could not delete a record SWIFT code NONEXISTENT not found",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("DELETE", ts.URL+"/v1/swift-codes/"+tc.input, nil)
			assert.NoError(t, err)

			resp, err := http.DefaultClient.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			var response map[string]string
			err = json.NewDecoder(resp.Body).Decode(&response)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, response)
		})
	}
}
