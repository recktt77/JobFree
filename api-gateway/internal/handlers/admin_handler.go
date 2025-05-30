package handlers

import (
	"api-gateway/clients"
	"net/http"

	pb "github.com/recktt77/projectProto-definitions/gen/admin_service"
	"github.com/gin-gonic/gin"
)

func BanUser(c *gin.Context) {
	var req pb.BanUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := clients.GetAdminClient().BanUser(c, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func DeleteReview(c *gin.Context) {
	var req pb.DeleteReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := clients.GetAdminClient().DeleteReview(c, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func ModerateProject(c *gin.Context) {
	var req pb.ModerateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := clients.GetAdminClient().ModerateProject(c, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func GetPlatformStats(c *gin.Context) {
	res, err := clients.GetAdminClient().GetPlatformStats(c, &pb.GetStatsRequest{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}
