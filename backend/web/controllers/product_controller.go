package controllers

import (
	"FlashSaleGo/common"
	"FlashSaleGo/model"
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

	c.HTML(http.StatusOK, "all", gin.H{
		"productArray": productArray,
	})
}

func (p *ProductController) GetAdd(c *gin.Context) {
	c.HTML(http.StatusOK, "add", gin.H{})
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
	c.HTML(http.StatusOK, "manager", gin.H{
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

func (p *ProductController) PostUpdate(c *gin.Context) {
	product := &model.Product{}
	c.Request.ParseForm()
	dec := common.NewDecoder(&common.DecoderOptions{TagName: "iRaiden"})
	if err := dec.Decode(c.Request.Form, product); err != nil {
		c.Error(err)
		return
	}
	err := p.ProductService.UpdateProduct(product)
	if err != nil {
		c.Error(err)
		return
	}
	c.Redirect(http.StatusMovedPermanently, "/product/all")
}

func (p *ProductController) PostAdd(c *gin.Context) {
	product := &model.Product{}
	c.Request.ParseForm()
	dec := common.NewDecoder(&common.DecoderOptions{TagName: "iRaiden"})
	if err := dec.Decode(c.Request.Form, product); err != nil {
		c.Error(err)
		return
	}
	_, err := p.ProductService.InsertProduct(product)
	if err != nil {
		c.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/product/all")
}
