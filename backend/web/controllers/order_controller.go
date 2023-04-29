package controllers

import (
	"FlashSaleGo/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type OrderController struct {
	OrderService service.IOrderService
}

func (o *OrderController) GetOrder(c *gin.Context) {
	orderArray, err := o.OrderService.GetAllOrder()
	if err != nil {
		c.Error(err)
	}
	c.HTML(http.StatusOK, "order", gin.H{
		"orderArray": orderArray,
	})
}
