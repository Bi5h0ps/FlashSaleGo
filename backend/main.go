package main

import (
	"FlashSaleGo/backend/web/controllers"
	"FlashSaleGo/repository"
	"FlashSaleGo/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	ginServer := gin.Default()
	//logger middleware
	ginServer.Use(gin.Logger())
	gin.SetMode(gin.DebugMode)
	// Load all templates from the views directory and subdirectories
	ginServer.LoadHTMLGlob("backend/web/views/**/*.html")

	ginServer.Static("/assets", "./backend/web/assets")
	// Register a layout
	ginServer.StaticFile("/", "backend/web/views/shared/layout.html")

	ginServer.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "shared/error.html", gin.H{
			"message": "There's an error on the requested page",
		})
	})
	ginServer.Use(func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.HTML(http.StatusInternalServerError, "shared/error.html", gin.H{
					"message": "Internal server error",
				})
			}
		}()

		c.Next()
	})

	// Set up product repository and service
	productRepo := repository.NewProductRepository(nil)
	productService := service.NewProductService(productRepo)

	// Set up product controller and register with Gin
	product := ginServer.Group("/product")
	{
		productController := controllers.ProductController{ProductService: productService}
		product.GET("/all", productController.GetAll)
		product.GET("/add", productController.GetAdd)
		product.GET("/manager", productController.GetAdd)
	}

	ginServer.Run(":8080")
}
