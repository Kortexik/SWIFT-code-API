package models

type DetailsResponse struct {
	Address       string        `json:"address"`
	BankName      string        `json:"bankName"`
	CountryISO2   string        `json:"countryISO2"`
	CountryName   string        `json:"countryName"`
	IsHeadquarter bool          `json:"isHeadquarter"`
	SwiftCode     string        `json:"swiftCode"`
	Branches      []BankDetails `json:"branches"`
}

type BranchResponse struct {
	Address       string `json:"address"`
	BankName      string `json:"bankName"`
	CountryISO2   string `json:"countryISO2"`
	CountryName   string `json:"countryName"`
	IsHeadquarter bool   `json:"isHeadquarter"`
	SwiftCode     string `json:"swiftCode"`
}

type CountryResponse struct {
	CountryISO2 string        `json:"countryISO2"`
	CountryName string        `json:"countryName"`
	SwiftCodes  []BankDetails `json:"swiftCodes"`
}

type BankDetails struct {
	Address       string `json:"address"`
	BankName      string `json:"bankName"`
	CountryISO2   string `json:"countryISO2"`
	IsHeadquarter bool   `json:"isHeadquarter"`
	SwiftCode     string `json:"swiftCode"`
}
