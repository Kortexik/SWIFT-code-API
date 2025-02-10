package main

import (
	"RemitlyTask/src/database"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {

	fmt.Println("Database connected:", database.DB != nil)
	r := gin.Default()
	r.Run(":8080")
}
