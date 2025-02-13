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
	service *services.SwiftCodeService
}

func NewSwiftCodeHandler(db *gorm.DB) *SwiftCodeHandler {
	repo := repositories.NewSwiftCodeRepository(db)
	service := services.NewSwiftCodeService(repo)
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

	countryResponse, ok := response.(models.CountryResponse)
	if !ok {
		log.Println("Invalid response type")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Invalid response type"})
		return
	}

	if countryResponse.CountryName == "" {
		log.Println("ISO2 code: " + iso2 + " is not valid.")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "ISO2 code " + iso2 + " is not valid."})
		return
	}

	c.JSON(http.StatusOK, countryResponse)
}

func (h *SwiftCodeHandler) AddNewSwiftCode(c *gin.Context) {
	var newSwiftCode models.SwiftCode

	if err := c.ShouldBindJSON(&newSwiftCode); err != nil {
		log.Println("Error binding JSON: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err := h.service.AddSwiftCode(&newSwiftCode)
	if err != nil {
		log.Println("Error inserting new code: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": ErrFailedToInsert + newSwiftCode.SwiftCode})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": newSwiftCode.SwiftCode + " has been added to database."})
}

func (h *SwiftCodeHandler) DeleteCode(c *gin.Context) {
	swiftCode := c.Param("swift-code")

	err := h.service.DeleteSwiftCode(swiftCode)
	if err != nil {
		log.Printf("Error deleting code %s: %v", swiftCode, err)
		c.JSON(http.StatusBadRequest, gin.H{"message": ErrFailedToDelete + ", given: " + swiftCode})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": swiftCode + " was removed."})
}
