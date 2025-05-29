// internal/middleware/validation.go
package middleware

import (
	"context"
	"errors"

	pb "github.com/recktt77/projectProto-definitions/gen/project_service/genproto/project"

	"google.golang.org/grpc"
)

func ValidationInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		switch r := req.(type) {
		case *pb.CreateProjectRequest:
			if r.ClientId == "" || r.Title == "" || r.Description == "" || len(r.Tasks) == 0 {
				return nil, errors.New("invalid CreateProjectRequest: required fields missing")
			}

		case *pb.UpdateProjectRequest:
			if r.ProjectId == "" || r.Title == "" || r.Description == "" {
				return nil, errors.New("invalid UpdateProjectRequest: required fields missing")
			}

		case *pb.DeleteProjectRequest:
			if r.Id == "" {
				return nil, errors.New("invalid DeleteProjectRequest: id is required")
			}

		case *pb.AttachFileRequest:
			if r.ProjectId == "" || r.FileUrl == "" {
				return nil, errors.New("invalid AttachFileRequest: project_id and file_url are required")
			}

		case *pb.DeleteFileRequest:
			if r.ProjectId == "" || r.FileUrl == "" {
				return nil, errors.New("invalid DeleteFileRequest: project_id and file_url are required")
			}

		case *pb.CreateTaskRequest:
			if r.ProjectId == "" || r.Title == "" || r.Description == "" {
				return nil, errors.New("invalid CreateTaskRequest: project_id, title, description are required")
			}

		case *pb.UpdateTaskRequest:
			if r.ProjectId == "" || r.TaskTitle == "" || r.NewTitle == "" {
				return nil, errors.New("invalid UpdateTaskRequest: required fields missing")
			}

		case *pb.DeleteTaskRequest:
			if r.ProjectId == "" || r.TaskTitle == "" {
				return nil, errors.New("invalid DeleteTaskRequest: required fields missing")
			}

		case *pb.GetTasksRequest:
			if r.ProjectId == "" {
				return nil, errors.New("invalid GetTasksRequest: project_id is required")
			}

		case *pb.GetProjectRequest:
			if r.Id == "" {
				return nil, errors.New("invalid GetProjectRequest: id is required")
			}
		}

		return handler(ctx, req)
	}
}
