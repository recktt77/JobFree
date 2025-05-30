// handlers/payment_handler.go
package handlers

import (
	"api-gateway/clients"
	"net/http"

	pb "github.com/recktt77/projectProto-definitions/gen/auth_service/genproto/payment"
	"github.com/gin-gonic/gin"
)

func CreatePayment(c *gin.Context) {
	var req pb.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := clients.GetPaymentClient().CreatePayment(c, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func GetPayment(c *gin.Context) {
	var req pb.GetPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := clients.GetPaymentClient().GetPayment(c, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func ListUserPayments(c *gin.Context) {
	var req pb.ListUserPaymentsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := clients.GetPaymentClient().ListUserPayments(c, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}
