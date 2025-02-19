package handlers

import (
	"RemitlyTask/src/models"
	"RemitlyTask/src/repositories"
	"RemitlyTask/src/services"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SwiftCodeHandler struct {
	service services.ISwiftCodeService
}

func NewSwiftCodeHandler(db *gorm.DB) *SwiftCodeHandler {
	repo := repositories.NewSwiftCodeRepository(db)
	service := services.NewSwiftCodeService(repo)
	return &SwiftCodeHandler{service: service}
}

func NewSwiftCodeHandlerByService(service services.ISwiftCodeService) *SwiftCodeHandler {
	return &SwiftCodeHandler{service: service}
}

func (h *SwiftCodeHandler) GetCode(c *gin.Context) {
	swiftCodeParam := c.Param("swift-code")

	if valid, response := validateSwiftCode(swiftCodeParam); !valid {
		c.JSON(http.StatusBadRequest, response)
		return

	}

	swiftCodePrefix, swiftCodeSuffix := parseSwiftCode(swiftCodeParam)
	if swiftCodeSuffix == "XXX" {
		response, err := h.service.GetHeadquarterDetails(swiftCodePrefix)
		if err != nil {
			log.Println(ErrFetchSwiftCodes, "for: ", swiftCodeParam, err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": ErrFetchSwiftCodes + "for: " + swiftCodeParam})
			return
		}

		if response == nil {
			c.JSON(http.StatusNotFound, gin.H{"message": ErrNoSwiftCodeFound + "for: " + swiftCodeParam})
			return
		}

		c.JSON(http.StatusOK, response)

	} else {
		response, err := h.service.GetBranchDetails(swiftCodeParam)
		if err != nil {
			log.Println(ErrFetchSwiftCodes, "for: ", swiftCodeParam, err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": ErrFetchSwiftCodes + "for: " + swiftCodeParam})
			return
		}

		if response == nil {
			c.JSON(http.StatusNotFound, gin.H{"message": ErrNoSwiftCodeFound + "for: " + swiftCodeParam})
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func (h *SwiftCodeHandler) GetCodesByCountry(c *gin.Context) {
	iso2 := strings.ToUpper(c.Param("ISO2"))

	if valid, response := validateISO2(iso2); !valid {
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response, err := h.service.GetSwiftCodesByCountry(iso2)
	if err != nil {
		log.Println("Error fetching swift codes:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": ErrFetchSwiftCodes + "for ISO2 code: " + iso2})
		return
	}

	SwiftCodeCountry, ok := response.(models.SwiftCodeCountry)
	if !ok {
		log.Println("Invalid response type")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Invalid response type"})
		return
	}

	if SwiftCodeCountry.CountryName == "" {
		log.Println("ISO2 code: " + iso2 + " is not valid.")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "ISO2 code " + iso2 + " is not valid."})
		return
	}

	c.JSON(http.StatusOK, SwiftCodeCountry)
}

func (h *SwiftCodeHandler) AddNewSwiftCode(c *gin.Context) {
	var newSwiftCode models.SwiftCodeBranch

	if err := c.ShouldBindJSON(&newSwiftCode); err != nil {
		log.Println("Error binding JSON: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if len(newSwiftCode.CountryISO2) != 2 {
		log.Println("Error inserting new code: this iso2 code does not exist.")
		c.JSON(http.StatusBadRequest, gin.H{"message": ErrFailedToInsert + "this iso2 code: " + newSwiftCode.CountryISO2 + " does not exist."})
		return
	}

	if len(newSwiftCode.SwiftCode) != 8 && len(newSwiftCode.SwiftCode) != 11 {
		log.Println("Error inserting new code: invalid swift code length")
		c.JSON(http.StatusBadRequest, gin.H{"message": ErrFailedToInsert + newSwiftCode.SwiftCode + " swift code must be 8 or 11 characters long."})
		return
	}

	if (!strings.HasSuffix(newSwiftCode.SwiftCode, "XXX") && newSwiftCode.IsHeadquarter) || (strings.HasSuffix(newSwiftCode.SwiftCode, "XXX") && !newSwiftCode.IsHeadquarter) {
		log.Println("Error inserting new code: isHeadquarter does not match the suffix of swiftcode.")
		c.JSON(http.StatusBadRequest, gin.H{"message": ErrFailedToInsert + "isHeadquarter does not match the suffix of swiftcode."})
		return
	}

	if newSwiftCode.Address == "" {
		log.Println("Error inserting new code: address can't be empty")
		c.JSON(http.StatusBadRequest, gin.H{"message": ErrFailedToInsert + "address can't be empty."})
		return
	}

	countryName, err := h.service.GetCountryName(newSwiftCode.CountryISO2)
	if err != nil {
		log.Println("Error checking country name from iso2")
		c.JSON(http.StatusBadRequest, gin.H{"message": ErrFailedToInsert + "invalid ISO2 code."})
		return
	}

	if !strings.EqualFold(newSwiftCode.CountryName, countryName) && countryName != "" {
		log.Println("Error inserting new code: iso2 code must match with given country.")
		c.JSON(http.StatusBadRequest, gin.H{"message": ErrFailedToInsert + "iso2 code must match with given country."})
		return
	}

	existingCode, err := h.service.GetBranchDetails(newSwiftCode.SwiftCode)
	if err == nil && existingCode != nil {
		log.Println("Error inserting new code: swift code already exists")
		c.JSON(http.StatusBadRequest, gin.H{"message": ErrFailedToInsert + "swift code already exists"})
		return
	}

	newValidatedCode := models.SwiftCode{
		Address:     newSwiftCode.Address,
		Name:        newSwiftCode.BankName,
		CountryISO2: newSwiftCode.CountryISO2,
		SwiftCode:   newSwiftCode.SwiftCode,
		CountryName: newSwiftCode.CountryName,
	}

	err = h.service.AddSwiftCode(&newValidatedCode)
	if err != nil {
		log.Println("Error inserting new code: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": ErrFailedToInsert + newSwiftCode.SwiftCode})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": newSwiftCode.SwiftCode + " has been added to the database."})
}

func (h *SwiftCodeHandler) DeleteCode(c *gin.Context) {
	swiftCode := c.Param("swift-code")

	err := h.service.DeleteSwiftCode(swiftCode)
	if err != nil {
		log.Printf("Error deleting code %s: %v", swiftCode, err)
		c.JSON(http.StatusNotFound, gin.H{"message": ErrFailedToDelete + " " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": swiftCode + " was removed."})
}
