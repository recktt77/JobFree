// internal/handler/project_handler.go
package handler

import (
	"context"
	"encoding/json"
	"log"
	"project_service/internal/cache"
	"project_service/internal/events"
	"project_service/internal/model"
	"project_service/internal/repository"
	"time"

	pb "github.com/recktt77/projectProto-definitions/gen/project_service/genproto/project"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProjectService struct {
	pb.UnimplementedProjectServiceServer
	repo      *repository.ProjectRepository
	publisher *events.Publisher
	cache     *cache.Redis
}

func New(repo *repository.ProjectRepository, publisher *events.Publisher, cache *cache.Redis) *ProjectService {
	return &ProjectService{repo: repo, publisher: publisher, cache: cache}
}

func (s *ProjectService) CreateProject(ctx context.Context, req *pb.CreateProjectRequest) (*pb.CreateProjectResponse, error) {
	if req.GetClientId() == "" || req.GetTitle() == "" || req.GetDescription() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "missing required fields")
	}
	tasks := make([]model.Task, len(req.Tasks))
	for i, t := range req.Tasks {
		tasks[i] = model.Task{Title: t.Title, Description: t.Description, Status: t.Status}
	}

	project := &model.Project{
		ClientID:    req.ClientId,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Tasks:       tasks,
		Attachments: req.Attachments,
	}

	id, err := s.repo.Create(ctx, project)
	if err != nil {
		return nil, err
	}
	project.ID = id
	_ = s.publisher.PublishProjectCreated(project)
	if err := s.publisher.PublishProjectCreated(project); err != nil {
		log.Println("failed to publish project.created:", err)
	}

	return &pb.CreateProjectResponse{Id: id, Message: "Project created successfully"}, nil
}

func (s *ProjectService) GetProject(ctx context.Context, req *pb.GetProjectRequest) (*pb.GetProjectResponse, error) {
	if req.GetId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "missing required fields")
	}

	cacheKey := "project:" + req.Id

	// Redis
	if cached, err := s.cache.Get(ctx, cacheKey); err == nil {
		var res pb.GetProjectResponse
		if json.Unmarshal([]byte(cached), &res) == nil {
			log.Println("GetProject from cache")
			return &res, nil
		}
	}

	// Mongo
	project, err := s.repo.FindByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	tasks := make([]*pb.Task, len(project.Tasks))
	for i, t := range project.Tasks {
		tasks[i] = &pb.Task{Title: t.Title, Description: t.Description, Status: t.Status}
	}
	resp := &pb.GetProjectResponse{
		Project: &pb.ProjectData{
			Id:          project.ID,
			ClientId:    project.ClientID,
			Title:       project.Title,
			Description: project.Description,
			Status:      project.Status,
			Attachments: project.Attachments,
			Tasks:       tasks,
		},
	}

	// Caching
	if data, err := json.Marshal(resp); err == nil {
		s.cache.Set(ctx, cacheKey, string(data), 10*time.Minute)
	}

	return resp, nil
}

func (s *ProjectService) GetAllProjects(ctx context.Context, req *pb.GetAllProjectsRequest) (*pb.GetAllProjectsResponse, error) {
	if req.GetClientId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "missing required fields")
	}

	cacheKey := "projects:client:" + req.ClientId

	// Redis
	if cached, err := s.cache.Get(ctx, cacheKey); err == nil {
		var res pb.GetAllProjectsResponse
		if json.Unmarshal([]byte(cached), &res) == nil {
			log.Println("GetAllProjects from cache")
			return &res, nil
		}
	}

	// Mongo
	projects, err := s.repo.GetAll(ctx, req.ClientId)
	if err != nil {
		return nil, err
	}

	resp := &pb.GetAllProjectsResponse{}
	for _, p := range projects {
		tasks := make([]*pb.Task, len(p.Tasks))
		for i, t := range p.Tasks {
			tasks[i] = &pb.Task{Title: t.Title, Description: t.Description, Status: t.Status}
		}
		resp.Projects = append(resp.Projects, &pb.ProjectData{
			Id:          p.ID,
			ClientId:    p.ClientID,
			Title:       p.Title,
			Description: p.Description,
			Status:      p.Status,
			Attachments: p.Attachments,
			Tasks:       tasks,
		})
	}

	// Caching
	if data, err := json.Marshal(resp); err == nil {
		s.cache.Set(ctx, cacheKey, string(data), 10*time.Minute)
	}

	return resp, nil
}

func (s *ProjectService) UpdateProject(ctx context.Context, req *pb.UpdateProjectRequest) (*pb.UpdateProjectResponse, error) {
	if req.GetClientId() == "" || req.GetProjectId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "missing required fields")
	}
	tasks := make([]model.Task, len(req.Tasks))
	for i, t := range req.Tasks {
		tasks[i] = model.Task{Title: t.Title, Description: t.Description, Status: t.Status}
	}

	project := &model.Project{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Tasks:       tasks,
		Attachments: req.Attachments,
	}

	err := s.repo.Update(ctx, req.ProjectId, project)
	if err != nil {
		return nil, err
	}
	return &pb.UpdateProjectResponse{Message: "Project updated successfully"}, nil
}

func (s *ProjectService) DeleteProject(ctx context.Context, req *pb.DeleteProjectRequest) (*pb.DeleteProjectResponse, error) {
	if req.GetId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "missing required fields")
	}
	err := s.repo.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteProjectResponse{Message: "Project deleted successfully"}, nil
}

func (s *ProjectService) AttachFile(ctx context.Context, req *pb.AttachFileRequest) (*pb.AttachFileResponse, error) {
	if req.GetProjectId() == "" || req.GetFileUrl() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "missing required fields")
	}
	err := s.repo.AttachFile(ctx, req.ProjectId, req.FileUrl)
	if err != nil {
		return nil, err
	}
	return &pb.AttachFileResponse{Message: "File attached successfully"}, nil
}

func (s *ProjectService) DeleteFile(ctx context.Context, req *pb.DeleteFileRequest) (*pb.DeleteFileResponse, error) {
	if req.GetProjectId() == "" || req.GetFileUrl() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "missing required fields")
	}
	err := s.repo.DeleteFile(ctx, req.ProjectId, req.FileUrl)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteFileResponse{Message: "File deleted successfully"}, nil
}

func (s *ProjectService) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
	if req.GetProjectId() == "" || req.GetTitle() == "" || req.GetDescription() == "" || req.GetStatus() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "missing required fields")
	}

	exists, erro := s.repo.HasTaskWithTitle(ctx, req.ProjectId, req.Title)
	if erro != nil {
		return nil, status.Errorf(codes.Internal, "failed to check task title: %v", erro)
	}
	if exists {
		return nil, status.Errorf(codes.AlreadyExists, "task title already exists in this project")
	}

	task := model.Task{Title: req.Title, Description: req.Description, Status: req.Status}
	err := s.repo.AddTask(ctx, req.ProjectId, task)
	if err != nil {
		return nil, err
	}
	return &pb.CreateTaskResponse{Message: "Task added successfully"}, nil
}

func (s *ProjectService) UpdateTask(ctx context.Context, req *pb.UpdateTaskRequest) (*pb.UpdateTaskResponse, error) {
	if req.GetProjectId() == "" || req.GetNewDescription() == "" || req.GetNewStatus() == "" || req.GetTaskTitle() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "missing required fields")
	}
	task := model.Task{Title: req.NewTitle, Description: req.NewDescription, Status: req.NewStatus}
	err := s.repo.UpdateTask(ctx, req.ProjectId, req.TaskTitle, task)
	if err != nil {
		return nil, err
	}
	return &pb.UpdateTaskResponse{Message: "Task updated successfully"}, nil
}

func (s *ProjectService) DeleteTask(ctx context.Context, req *pb.DeleteTaskRequest) (*pb.DeleteTaskResponse, error) {
	if req.GetTaskTitle() == "" || req.GetProjectId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "missing required fields")
	}
	err := s.repo.DeleteTask(ctx, req.ProjectId, req.TaskTitle)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteTaskResponse{Message: "Task deleted successfully"}, nil
}

func (s *ProjectService) GetTasks(ctx context.Context, req *pb.GetTasksRequest) (*pb.GetTasksResponse, error) {
	if req.GetProjectId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "missing required fields")
	}

	cacheKey := "tasks:project:" + req.ProjectId

	// Redis
	if cached, err := s.cache.Get(ctx, cacheKey); err == nil {
		var res pb.GetTasksResponse
		if json.Unmarshal([]byte(cached), &res) == nil {
			log.Println("GetTasks from cache")
			return &res, nil
		}
	}

	// Mongo
	tasksData, err := s.repo.GetTasks(ctx, req.ProjectId)
	if err != nil {
		return nil, err
	}

	tasks := make([]*pb.Task, len(tasksData))
	for i, t := range tasksData {
		tasks[i] = &pb.Task{Title: t.Title, Description: t.Description, Status: t.Status}
	}

	resp := &pb.GetTasksResponse{Tasks: tasks}

	// Caching
	if data, err := json.Marshal(resp); err == nil {
		s.cache.Set(ctx, cacheKey, string(data), 10*time.Minute)
	}

	return resp, nil
}
