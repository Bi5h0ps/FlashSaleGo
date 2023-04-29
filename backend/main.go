package main

import (
	"FlashSaleGo/backend/web/controllers"
	"FlashSaleGo/repository"
	"FlashSaleGo/service"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	ginServer := gin.Default()
	//logger middleware
	ginServer.Use(gin.Logger())
	gin.SetMode(gin.DebugMode)

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

	ginServer.LoadHTMLGlob("backend/web/views/**/*.html")
	ginServer.Static("/assets", "./backend/web/assets")

	// Define multitemplate templates
	templates := multitemplate.New()
	// Add the parent template
	templates.AddFromFiles("add", "backend/web/views/shared/layout.html",
		"backend/web/views/product/add.html")
	templates.AddFromFiles("all", "backend/web/views/shared/layout.html",
		"backend/web/views/product/view.html")
	templates.AddFromFiles("manager", "backend/web/views/shared/layout.html",
		"backend/web/views/product/manager.html")
	templates.AddFromFiles("order", "backend/web/views/shared/layout.html",
		"backend/web/views/order/view.html")
	// Set the router HTML templates to the multitemplate engine
	ginServer.HTMLRender = templates

	repoProduct := repository.NewProductRepository(nil)
	serviceProduct := service.NewProductService(repoProduct)
	controllerProduct := controllers.ProductController{ProductService: serviceProduct}

	product := ginServer.Group("product")
	{
		product.GET("/all", controllerProduct.GetAll)
		product.GET("/add", controllerProduct.GetAdd)
		product.GET("/manager", controllerProduct.GetManager)
		product.POST("/update", controllerProduct.PostUpdate)
		product.POST("/add", controllerProduct.PostAdd)
		product.GET("/delete", controllerProduct.GetDelete)
	}

	repoOrder := repository.NewOrderRepository(nil)
	serviceOrder := service.NewOrderService(repoOrder)
	controllerOrder := controllers.OrderController{OrderService: serviceOrder}

	order := ginServer.Group("order")
	{
		order.GET("/", controllerOrder.GetOrder)
	}
	ginServer.Run(":8080")
}
