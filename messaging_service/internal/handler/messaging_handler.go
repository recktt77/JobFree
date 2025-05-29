package handler

import (
	"context"

	"messaging_service/internal/model"
	"messaging_service/internal/repository"

	"github.com/recktt77/projectProto-definitions/gen/messaging_service/genproto/messaging"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessagingHandler struct {
	messaging.UnimplementedMessagingServiceServer
	repo repository.MessageRepository
}

func NewMessagingHandler(repo repository.MessageRepository) *MessagingHandler {
	return &MessagingHandler{
		repo: repo,
	}
}

func (h *MessagingHandler) SendMessage(ctx context.Context, req *messaging.SendMessageRequest) (*messaging.SendMessageResponse, error) {
	conversationID, err := primitive.ObjectIDFromHex(req.GetConversationId())
	if err != nil {
		return nil, err
	}
	senderID, err := primitive.ObjectIDFromHex(req.GetSenderId())
	if err != nil {
		return nil, err
	}

	msg := &model.Message{
		ConversationID: conversationID,
		SenderID:       senderID,
		Content:        req.GetContent(),
	}

	id, err := h.repo.CreateMessage(ctx, msg)
	if err != nil {
		return nil, err
	}

	return &messaging.SendMessageResponse{
		MessageId: id.Hex(),
		Message:   "Message sent successfully",
	}, nil
}

func (h *MessagingHandler) GetMessages(ctx context.Context, req *messaging.GetMessagesRequest) (*messaging.GetMessagesResponse, error) {
	conversationID, err := primitive.ObjectIDFromHex(req.GetConversationId())
	if err != nil {
		return nil, err
	}

	msgs, err := h.repo.GetMessagesByConversation(ctx, conversationID)
	if err != nil {
		return nil, err
	}

	var responseMessages []*messaging.Message
	for _, m := range msgs {
		responseMessages = append(responseMessages, &messaging.Message{
			Id:             m.ID.Hex(),
			ConversationId: m.ConversationID.Hex(),
			SenderId:       m.SenderID.Hex(),
			Content:        m.Content,
			SentAt:         m.SentAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	return &messaging.GetMessagesResponse{
		Messages: responseMessages,
	}, nil
}
