package handlers

import (
	"api-gateway/clients"
	"net/http"

	matching "github.com/recktt77/projectProto-definitions/gen/matching_service/recktt77/projectProto-definitions/matching_service"

	"github.com/gin-gonic/gin"
)

func CreateBid(c *gin.Context) {
	var req matching.CreateBidRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := clients.GetMatchingClient().CreateBid(c, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func GetBidsForProject(c *gin.Context) {
	projectID := c.Param("project_id")

	res, err := clients.GetMatchingClient().GetBidsForProject(c, &matching.GetBidsRequest{
		ProjectId: projectID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func MatchFreelancers(c *gin.Context) {
	var req matching.MatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := clients.GetMatchingClient().MatchFreelancers(c, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
