package routers

import (
	"assignment-2/controllers"

	"github.com/gin-gonic/gin"
)

func StartServer() *gin.Engine {
	router := gin.Default()

	// Create
	router.POST("/orders", controllers.CreateOrder)
	// Update
	router.PUT("/orders/:orderID", controllers.UpdateOrder)
	// Read
	router.GET("/orders/:orderID", controllers.GetOrder)
	// Delete
	router.DELETE("/orders/:orderID", controllers.DeleteOrder)
	// Read All
	router.GET("/orders", controllers.GetAllOrder)

	return router
}
