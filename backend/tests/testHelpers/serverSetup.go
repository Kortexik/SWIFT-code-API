package testHelpers

import (
	"RemitlyTask/src/handlers"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

func SetupTestServer(handler *handlers.SwiftCodeHandler) *httptest.Server {
	r := gin.Default()
	vCodes := r.Group("v1/swift-codes")
	{
		vCodes.POST("/", handler.AddNewSwiftCode)
		vCodes.GET("/:swift-code", handler.GetCode)
		vCodes.GET("/country/:ISO2", handler.GetCodesByCountry)
		vCodes.DELETE("/:swift-code", handler.DeleteCode)
	}
	return httptest.NewServer(r)
}
