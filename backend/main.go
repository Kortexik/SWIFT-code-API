package main

import (
	"RemitlyTask/src/database"
	"RemitlyTask/src/handlers"

	"github.com/gin-gonic/gin"
)

func main() {

	handler := handlers.NewSwiftCodeHandler(database.DB)
	r := gin.Default()

	vCodes := r.Group("v1/swift-codes")
	{
		vCodes.GET("/:swift-code", handler.GetCode)
		vCodes.GET("/country/:ISO2", handler.GetCodesByCountry)
		vCodes.POST("", handler.AddNewSwiftCode)
		vCodes.DELETE("/:swift-code", handler.DeleteCode)
	}
	r.Run(":8080")
}
