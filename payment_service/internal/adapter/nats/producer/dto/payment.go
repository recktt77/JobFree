package dto

type PaymentCreated struct {
	ID             string  `json:"id"`
	UserID         string  `json:"user_id"`
	SubscriptionID string  `json:"subscription_id"`
	Amount         float64 `json:"amount"`
	Status         string  `json:"status"`
	CreatedAt      string  `json:"created_at"`
}
