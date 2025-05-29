package models

type Bid struct {
	BidID        string `bson:"bidid"`
	ProjectID    string `bson:"projectid"`
	FreelancerID string `bson:"freelancerid"`
	CoverLetter  string `bson:"coverletter"`
	Timestamp    string `bson:"timestamp"`
}

type Freelancer struct {
	FreelancerID string   `bson:"freelancerid"`
	Name         string   `bson:"name"`
	Skills       []string `bson:"skills"`
}
