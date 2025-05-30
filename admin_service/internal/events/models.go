package events

type UserBannedEvent struct {
	UserID    string `json:"user_id"`
	Reason    string `json:"reason"`
	Timestamp string `json:"timestamp"`
}
