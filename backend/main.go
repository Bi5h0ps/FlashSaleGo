package main

import (
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
	templates.AddFromFiles("order", "backend/web/views/shared/layout.html",
		"backend/web/views/order/view.html")
	// Set the router HTML templates to the multitemplate engine
	ginServer.HTMLRender = templates
	ginServer.GET("/product/add", func(context *gin.Context) {
		// Render the view.html child template inside the layout.html parent template
		context.HTML(http.StatusOK, "add", gin.H{})
	})
	ginServer.GET("/product/all", func(context *gin.Context) {
		// Render the view.html child template inside the layout.html parent template
		context.HTML(http.StatusOK, "all", gin.H{})
	})
	ginServer.GET("/order", func(context *gin.Context) {
		// Render the view.html child template inside the layout.html parent template
		context.HTML(http.StatusOK, "order", gin.H{})
	})
	ginServer.Run(":8080")
}
