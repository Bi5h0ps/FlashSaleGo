package controllers

import (
	"FlashSaleGo/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ProductController struct {
	ProductService service.IProductService
}

func (p *ProductController) GetAll(c *gin.Context) {
	productArray, err := p.ProductService.GetAllProduct()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get products",
		})
		return
	}

	c.HTML(http.StatusOK, "product/view.html", gin.H{
		"productArray": productArray,
	})
}

func (p *ProductController) GetAdd(c *gin.Context) {
	c.HTML(http.StatusOK, "product/add.html", gin.H{})
}

func (p *ProductController) GetManager(c *gin.Context) {
	idstring := c.Query("id")
	id, err := strconv.ParseInt(idstring, 10, 16)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	product, err := p.ProductService.GetProductById(id)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.HTML(http.StatusOK, "product/manager.tmpl", gin.H{
		"product": product,
	})
}

func (p *ProductController) GetDelete(ctx *gin.Context) {
	idstring := ctx.Query("id")
	id, err := strconv.ParseInt(idstring, 10, 16)
	if err != nil {
		ctx.Error(err)
		return
	}
	p.ProductService.DeleteProductById(id)
	ctx.Redirect(http.StatusFound, "/product/all")
}
