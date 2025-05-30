package dao

import (
	"context"
	"errors"

	"github.com/recktt77/JobFree/payment_service/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type paymentRepository struct {
	collection *mongo.Collection
}

func NewPaymentRepository(db *mongo.Database) *paymentRepository {
	return &paymentRepository{
		collection: db.Collection("payments"),
	}
}

func (r *paymentRepository) CreatePayment(ctx context.Context, payment model.Payment) error {
	_, err := r.collection.InsertOne(ctx, payment)
	return err
}

func (r *paymentRepository) GetPayment(ctx context.Context, id string) (*model.Payment, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var payment model.Payment
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&payment)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return &payment, nil
}

func (r *paymentRepository) ListUserPayments(ctx context.Context, userID string) ([]model.Payment, error) {
	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	cursor, err := r.collection.Find(ctx, bson.M{"user_id": uid})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var payments []model.Payment
	if err := cursor.All(ctx, &payments); err != nil {
		return nil, err
	}

	return payments, nil
}
