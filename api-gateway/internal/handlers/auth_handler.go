package handlers

import (
	"api-gateway/clients"
	"net/http"

	"github.com/recktt77/projectProto-definitions/gen/auth_service/genproto/auth"

	"github.com/gin-gonic/gin"
)

func RegisterUser(c *gin.Context) {
	var req auth.RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := clients.GetAuthClient().RegisterUser(c, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func LoginUser(c *gin.Context) {
	var req auth.LoginUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := clients.GetAuthClient().LoginUser(c, &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
