package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rohyunge/main/api"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	router := gin.Default()

	router.Use(CORSMiddleware())

	router.GET("/reserve/:reserveId", api.GetReservation)
	router.GET("/reserve", api.GetReservationList)
	router.POST("/reserve", api.AddReservation)
	router.PUT("/reserve/:reserveId", api.ModifyReservation)
	router.DELETE("/reserve/:reserveId", api.DeleteReservation)

	router.Run()

}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, Origin")
		// c.Header("Access-Control-Allow-Credentials", "false") // false에  Access-Control-Allow-Origin을 "*"로주면 모든 URL 허용
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, DELETE, POST, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
