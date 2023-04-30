package controller

import (
	"FlashSaleGo/grpc/order"
	"FlashSaleGo/model"
	"FlashSaleGo/service"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"strconv"
)

type ProductController struct {
	ProductService service.IProductService
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

func (p *ProductController) GetOrder(c *gin.Context) {
	var userName string
	if user, exists := c.Get("user"); exists {
		if u, ok := user.(*model.User); ok {
			userName = u.UserName
		}
	}
	productID, err := strconv.ParseInt(c.Query("productID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "productID provided in wrong format",
		})
		return
	}
	//rpc call
	var conn *grpc.ClientConn
	conn, err = grpc.Dial(":9093", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()
	client := order.NewOrderServiceClient(conn)
	response, err := client.MakeOrder(context.Background(), &order.OrderInfo{Username: userName, ProductID: productID})
	if err != nil {
		log.Fatalf("Error when calling MakeOrder: %s", err)
	}
	log.Printf("Order Status from server: %s", response.IsOrderSuccess)

	c.HTML(http.StatusOK, "result", gin.H{
		"showMessage": response.IsOrderSuccess,
		"orderID":     "000",
	})
}
