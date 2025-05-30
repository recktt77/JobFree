package handlers

import (
	"api-gateway/clients"
	"net/http"

	"github.com/recktt77/projectProto-definitions/gen/review_service/genproto/review"

	"github.com/gin-gonic/gin"
)

func LeaveReview(c *gin.Context) {
	var req review.LeaveReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := clients.GetReviewClient().LeaveReview(c, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func GetReviews(c *gin.Context) {
	req := &review.GetReviewsRequest{
		ProjectId:  c.Query("project_id"),
		RevieweeId: c.Query("reviewee_id"),
	}

	res, err := clients.GetReviewClient().GetReviews(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func ModerateReview(c *gin.Context) {
	var req review.ModerateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := clients.GetReviewClient().ModerateReview(c, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
