package handlers

import "github.com/gin-gonic/gin"

func parseSwiftCode(swiftCode string) (prefix, suffix string) {
	return swiftCode[:len(swiftCode)-3], swiftCode[len(swiftCode)-3:]
}

func isValidSwiftCode(swiftCode string) bool {
	return len(swiftCode) == 8 || len(swiftCode) == 11
}

func validateSwiftCode(swiftCode string) (bool, *gin.H) {
	if !isValidSwiftCode(swiftCode) {
		return false, &gin.H{"message": ErrInvalidSwiftLength}
	}
	return true, nil
}

func isISO2Valid(iso2Code string) bool {
	return len(iso2Code) == 2
}

func validateISO2(iso2Code string) (bool, *gin.H) {
	if !isISO2Valid(iso2Code) {
		return false, &gin.H{"message": ErrInvalidISO2Length}
	}
	return true, nil
}
