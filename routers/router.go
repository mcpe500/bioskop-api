package routers

import (
	"bioskop-api/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// CORS Middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Grouping routes supaya lebih rapi
	api := r.Group("/bioskop")
	{
		api.POST("", controllers.CreateBioskop)       // Create
		api.GET("", controllers.GetAllBioskop)        // Read All
		api.GET("/:id", controllers.GetBioskopByID)   // Read One
		api.PUT("/:id", controllers.UpdateBioskop)    // Update
		api.DELETE("/:id", controllers.DeleteBioskop) // Delete
	}

	return r
}
