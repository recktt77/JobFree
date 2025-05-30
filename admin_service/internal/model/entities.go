package model

type User struct {
	ID     string `bson:"_id"`
	Banned bool   `bson:"banned"`
}

type Review struct {
	ID          string `bson:"_id"`
	Deleted     bool   `bson:"deleted"`
	ModeratorID string `bson:"moderator_id"`
}

type Project struct {
	ID     string `bson:"_id"`
	Status string `bson:"status"`
}

type PlatformStats struct {
	TotalUsers     int32
	BannedUsers    int32
	ActiveProjects int32
	TotalReviews   int32
}
