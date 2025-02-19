package repositories

import (
	"RemitlyTask/src/models"
	"fmt"

	"gorm.io/gorm"
)

type ISwiftCodeRepository interface {
	FindBySwiftCodePrefix(prefix string) ([]models.SwiftCode, error)
	FindBySwiftCode(code string) (models.SwiftCode, error)
	FindCountryNameByISO2(iso2 string) (string, error)
	FindByCountryISO2(iso2 string) ([]models.SwiftCode, error)
	Create(newCode *models.SwiftCode) error
	Delete(swiftCode string) error
}

type SwiftCodeRepository struct {
	db *gorm.DB
}

func NewSwiftCodeRepository(db *gorm.DB) ISwiftCodeRepository {
	return &SwiftCodeRepository{db: db}
}

func (r *SwiftCodeRepository) FindBySwiftCodePrefix(prefix string) ([]models.SwiftCode, error) {
	var swiftCodes []models.SwiftCode
	result := r.db.Where("swift_code LIKE ?", prefix+"%").Find(&swiftCodes)
	return swiftCodes, result.Error
}

func (r *SwiftCodeRepository) FindBySwiftCode(code string) (models.SwiftCode, error) {
	var swiftCode models.SwiftCode
	result := r.db.Where("swift_code = ?", code).Find(&swiftCode)
	return swiftCode, result.Error
}

func (r *SwiftCodeRepository) FindCountryNameByISO2(iso2 string) (string, error) {
	var countryName string
	result := r.db.Table("swift_codes").Select("country_name").Where("country_iso2 = ?", iso2).Scan(&countryName)
	return countryName, result.Error
}

func (r *SwiftCodeRepository) FindByCountryISO2(iso2 string) ([]models.SwiftCode, error) {
	var swiftCodes []models.SwiftCode
	result := r.db.Where("country_iso2 = ?", iso2).Find(&swiftCodes)
	return swiftCodes, result.Error
}

func (r *SwiftCodeRepository) Create(newCode *models.SwiftCode) error {
	return r.db.Create(newCode).Error
}

func (r *SwiftCodeRepository) Delete(swiftCode string) error {
	result := r.db.Where("swift_code = ?", swiftCode).Delete(&models.SwiftCode{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("SWIFT code %s not found", swiftCode)
	}
	return nil
}
