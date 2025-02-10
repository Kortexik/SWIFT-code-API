package models

type SwiftCode struct {
	ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	CountryISO2 string `gorm:"type:char(2);not null;check:length(country_iso2) = 2" json:"country_iso2"`
	SwiftCode   string `gorm:"type:varchar(11);not null;unique;check:length(swift_code) in (8, 11)" json:"swift_code"`
	CodeType    string `gorm:"type:varchar(5);not null" json:"code_type"`
	Name        string `gorm:"type:text;not null" json:"name"`
	Address     string `gorm:"type:text" json:"address"`
	TownName    string `gorm:"type:varchar(60)" json:"town_name"`
	CountryName string `gorm:"type:varchar(50)" json:"country_name"`
	TimeZone    string `gorm:"type:varchar(50)" json:"time_zone"`
}
