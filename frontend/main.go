package main

import (
	"FlashSaleGo/frontend/web/controller"
	"FlashSaleGo/initialization"
	"FlashSaleGo/middleware"
	"FlashSaleGo/repository"
	"FlashSaleGo/service"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	initialization.LoadEnvVariables()
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

	repoUser := repository.NewUserRepository(nil)
	serviceUser := service.NewUserService(repoUser)
	controllerUser := controller.UserController{UserService: serviceUser}

	product := ginServer.Group("product")
	product.Use(middleware.RequireAuth(serviceUser))
	{
		product.GET("/detail", controllerProduct.GetDetail)
	}

	user := ginServer.Group("user")
	{
		user.GET("/login", controllerUser.GetLogin)
		user.POST("/login", controllerUser.PostLogin)
		user.GET("/register", controllerUser.GetRegister)
		user.POST("/register", controllerUser.PostRegister)
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
