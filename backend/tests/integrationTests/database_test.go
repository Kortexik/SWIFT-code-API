package integrationTests

import (
	"RemitlyTask/src/models"
	"RemitlyTask/tests/testHelpers"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatabaseSetup(t *testing.T) {
	db := testHelpers.SetupTestDB(t)
	defer testHelpers.CleanupTestDB(t, db)

	sqlDB, err := db.DB()
	assert.NoError(t, err)
	err = sqlDB.Ping()
	assert.NoError(t, err, "Failed to ping the database")

	hasTable := db.Migrator().HasTable(&models.SwiftCode{})
	assert.True(t, hasTable, "SwiftCode table does not exist")

	testData := models.SwiftCode{
		Address:     "TEST ADDRESS",
		Name:        "TEST BANK",
		CountryISO2: "PL",
		CountryName: "POLAND",
		SwiftCode:   "TESTTESTTES",
	}

	err = db.Create(&testData).Error
	assert.NoError(t, err, "Failed to insert test data into the database")

	var retrievedData models.SwiftCode
	err = db.First(&retrievedData, "swift_code = ?", testData.SwiftCode).Error
	assert.NoError(t, err, "Failed to retrieve test data from the database")
	assert.Equal(t, testData, retrievedData, "Retrieved data does not match inserted data")
}
