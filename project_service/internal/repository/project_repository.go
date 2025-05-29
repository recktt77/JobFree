// internal/repository/project_repository.go
package repository

import (
	"context"
	"project_service/internal/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProjectRepository struct {
	collection *mongo.Collection
}

func NewProjectRepository(db *mongo.Database) *ProjectRepository {
	return &ProjectRepository{
		collection: db.Collection("projects"),
	}
}

func (r *ProjectRepository) HasTaskWithTitle(ctx context.Context, projectID string, title string) (bool, error) {
	filter := bson.M{
		"_id":         projectID,
		"tasks.title": title,
	}
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *ProjectRepository) Create(ctx context.Context, p *model.Project) (string, error) {
	p.ID = generateUUID()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_, err := r.collection.InsertOne(ctx, p)
	return p.ID, err
}

func (r *ProjectRepository) FindByID(ctx context.Context, id string) (*model.Project, error) {
	var result model.Project
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&result)
	return &result, err
}

func (r *ProjectRepository) GetAll(ctx context.Context, clientID string) ([]*model.Project, error) {
	filter := bson.M{}
	if clientID != "" {
		filter["client_id"] = clientID
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var projects []*model.Project
	for cursor.Next(ctx) {
		var p model.Project
		if err := cursor.Decode(&p); err == nil {
			projects = append(projects, &p)
		}
	}
	return projects, nil
}

func (r *ProjectRepository) Update(ctx context.Context, id string, p *model.Project) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_, err := r.collection.UpdateByID(ctx, id, bson.M{"$set": bson.M{
		"title": p.Title, "description": p.Description, "status": p.Status, "tasks": p.Tasks, "attachments": p.Attachments,
	}})
	return err
}

func (r *ProjectRepository) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *ProjectRepository) AttachFile(ctx context.Context, projectID, fileURL string) error {
	_, err := r.collection.UpdateByID(ctx, projectID, bson.M{"$push": bson.M{"attachments": fileURL}})
	return err
}

func (r *ProjectRepository) DeleteFile(ctx context.Context, projectID, fileURL string) error {
	_, err := r.collection.UpdateByID(ctx, projectID, bson.M{"$pull": bson.M{"attachments": fileURL}})
	return err
}

func (r *ProjectRepository) AddTask(ctx context.Context, projectID string, task model.Task) error {
	_, err := r.collection.UpdateByID(ctx, projectID, bson.M{"$push": bson.M{"tasks": task}})
	return err
}

func (r *ProjectRepository) UpdateTask(ctx context.Context, projectID, taskTitle string, updated model.Task) error {
	filter := bson.M{"_id": projectID, "tasks.title": taskTitle}
	update := bson.M{"$set": bson.M{"tasks.$.title": updated.Title, "tasks.$.description": updated.Description, "tasks.$.status": updated.Status}}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *ProjectRepository) DeleteTask(ctx context.Context, projectID, taskTitle string) error {
	_, err := r.collection.UpdateByID(ctx, projectID, bson.M{"$pull": bson.M{"tasks": bson.M{"title": taskTitle}}})
	return err
}

func (r *ProjectRepository) GetTasks(ctx context.Context, projectID string) ([]model.Task, error) {
	var project model.Project
	err := r.collection.FindOne(ctx, bson.M{"_id": projectID}).Decode(&project)
	return project.Tasks, err
}

func generateUUID() string {
	// stub: замените на github.com/google/uuid или другой генератор
	return time.Now().Format("20060102150405")
}
