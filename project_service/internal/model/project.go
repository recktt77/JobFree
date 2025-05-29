package model

type Task struct {
	Title       string `bson:"title"`
	Description string `bson:"description"`
	Status      string `bson:"status"`
}

type Project struct {
	ID          string   `bson:"_id,omitempty"`
	ClientID    string   `bson:"client_id"`
	Title       string   `bson:"title"`
	Description string   `bson:"description"`
	Status      string   `bson:"status"`
	Tasks       []Task   `bson:"tasks"`
	Attachments []string `bson:"attachments"`
}
