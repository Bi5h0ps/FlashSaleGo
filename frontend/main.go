package main

import (
	"FlashSaleGo/frontend/web/controller"
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
		c.HTML(http.StatusNotFound, "error", gin.H{
			"message": "Requested routing not exist",
		})
	})
	ginServer.Use(func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.HTML(http.StatusInternalServerError, "error", gin.H{
					"message": err,
				})
			}
		}()

		c.Next()
	})

	ginServer.LoadHTMLGlob("frontend/web/views/**/*.html")
	ginServer.Static("/public", "./frontend/web/public")

	// Set the router HTML templates to the multitemplate engine
	ginServer.HTMLRender = loadTemplates()

	repoProduct := repository.NewProductRepository(nil)
	serviceProduct := service.NewProductService(repoProduct)
	controllerProduct := controller.ProductController{ProductService: serviceProduct}

	product := ginServer.Group("product")
	{
		product.GET("/detail", controllerProduct.GetDetail)
	}

	ginServer.Run(":9092")
}

func loadTemplates() (templates multitemplate.Renderer) {
	// Define multitemplate templates
	templates = multitemplate.New()
	// Add the parent template
	templates.AddFromFiles("detail", "frontend/web/views/shared/productLayout.html",
		"frontend/web/views/product/view.html")
	templates.AddFromFiles("result", "frontend/web/views/shared/productLayout.html",
		"frontend/web/views/product/result.html")
	templates.AddFromFiles("login", "frontend/web/views/shared/layout.html",
		"frontend/web/views/user/login.html")
	templates.AddFromFiles("register", "frontend/web/views/shared/layout.html",
		"frontend/web/views/user/register.html")
	templates.AddFromFiles("error", "frontend/web/views/shared/error.html")
	return
}
