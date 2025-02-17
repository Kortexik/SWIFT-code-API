package services

import (
	"RemitlyTask/src/models"
	"RemitlyTask/src/repositories"
)

type SwiftCodeService struct {
	repo *repositories.SwiftCodeRepository
}

func NewSwiftCodeService(repo *repositories.SwiftCodeRepository) *SwiftCodeService {
	return &SwiftCodeService{repo: repo}
}

func (s *SwiftCodeService) GetHeadquarterDetails(swiftCodePrefix string) (interface{}, error) {
	swiftCodes, err := s.repo.FindBySwiftCodePrefix(swiftCodePrefix)
	if err != nil {
		return nil, err
	}

	if len(swiftCodes) == 0 {
		return nil, nil
	}

	var headquarter *models.SwiftCode
	var branches []models.SwiftCodeBank

	for _, code := range swiftCodes {
		if code.IsHeadquarter() {
			headquarter = &code
		} else {
			branches = append(branches, models.SwiftCodeBank{
				Address:       code.Address,
				BankName:      code.Name,
				CountryISO2:   code.CountryISO2,
				IsHeadquarter: false,
				SwiftCode:     code.SwiftCode,
			})
		}
	}

	if headquarter == nil {
		return nil, nil
	}

	detailedResponse := models.SwiftCodeDetails{
		Address:       headquarter.Address,
		BankName:      headquarter.Name,
		CountryISO2:   headquarter.CountryISO2,
		CountryName:   headquarter.CountryName,
		IsHeadquarter: true,
		SwiftCode:     headquarter.SwiftCode,
		Branches:      branches,
	}

	return detailedResponse, nil
}

func (s *SwiftCodeService) GetBranchDetails(swiftCode string) (interface{}, error) {
	branch, err := s.repo.FindBySwiftCode(swiftCode)
	if err != nil {
		return nil, err
	}

	if branch.SwiftCode == "" {
		return nil, nil
	}

	response := models.SwiftCodeBranch{
		Address:       branch.Address,
		BankName:      branch.Name,
		CountryISO2:   branch.CountryISO2,
		CountryName:   branch.CountryName,
		IsHeadquarter: false,
		SwiftCode:     branch.SwiftCode,
	}

	return response, nil
}

func (s *SwiftCodeService) GetSwiftCodesByCountry(iso2 string) (interface{}, error) {
	countryName, err := s.repo.FindCountryNameByISO2(iso2)
	if err != nil {
		return nil, err
	}

	swiftCodes, err := s.repo.FindByCountryISO2(iso2)
	if err != nil {
		return nil, err
	}

	var SwiftCodeBranchs []models.SwiftCodeBank

	for _, code := range swiftCodes {
		SwiftCodeBranchs = append(SwiftCodeBranchs, models.SwiftCodeBank{
			Address:       code.Address,
			BankName:      code.Name,
			CountryISO2:   code.CountryISO2,
			IsHeadquarter: code.IsHeadquarter(),
			SwiftCode:     code.SwiftCode,
		})
	}

	return models.SwiftCodeCountry{
		CountryISO2: iso2,
		CountryName: countryName,
		SwiftCodes:  SwiftCodeBranchs,
	}, nil
}

func (s *SwiftCodeService) AddSwiftCode(newCode *models.SwiftCode) error {
	return s.repo.Create(newCode)
}

func (s *SwiftCodeService) DeleteSwiftCode(swiftCode string) error {
	return s.repo.Delete(swiftCode)
}
