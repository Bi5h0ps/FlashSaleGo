package controller

import (
	"FlashSaleGo/rabbitmq"
	"FlashSaleGo/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ProductController struct {
	ProductService service.IProductService
	RabbitMQ       *rabbitmq.RabbitMQ
	OrderService   service.IOrderService
}

func (p *ProductController) GetDetail(c *gin.Context) {
	id, err := strconv.ParseInt(c.Query("productID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "productID provided in wrong format",
		})
		return
	}
	product, err := p.ProductService.GetProductById(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "productID incorrect, no such item found",
		})
		return
	}
	c.HTML(http.StatusOK, "detail", gin.H{
		"product": product,
	})
}
