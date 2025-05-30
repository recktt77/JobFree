package handlers

import (
	"api-gateway/clients"
	"net/http"

	"github.com/recktt77/projectProto-definitions/gen/messaging_service/genproto/messaging"

	"github.com/gin-gonic/gin"
)

func SendMessage(c *gin.Context) {
	var req messaging.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := clients.GetMessagingClient().SendMessage(c, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func GetMessages(c *gin.Context) {
	conversationID := c.Param("conversation_id")

	res, err := clients.GetMessagingClient().GetMessages(c, &messaging.GetMessagesRequest{
		ConversationId: conversationID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
