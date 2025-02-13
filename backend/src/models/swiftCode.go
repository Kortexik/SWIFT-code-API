package models

import "strings"

type SwiftCode struct {
	ID          uint   `gorm:"primaryKey;autoIncrement" json:"-"`
	Address     string `gorm:"type:text" json:"address"`
	Name        string `gorm:"type:text;not null" json:"bankName"`
	CountryISO2 string `gorm:"type:char(2);not null;check:length(country_iso2) = 2" json:"countryISO2"`
	SwiftCode   string `gorm:"type:varchar(11);not null;unique;check:length(swift_code) in (8, 11)" json:"swiftCode"`
	CodeType    string `gorm:"type:varchar(5);not null" json:"-"`
	TownName    string `gorm:"type:varchar(60)" json:"-"`
	CountryName string `gorm:"type:varchar(50)" json:"countryName,omitempty"`
	TimeZone    string `gorm:"type:varchar(50)" json:"-"`
}

func (s *SwiftCode) IsHeadquarter() bool {
	return strings.HasSuffix(s.SwiftCode, "XXX")
}
