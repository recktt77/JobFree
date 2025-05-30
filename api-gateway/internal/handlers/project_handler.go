package handlers

import (
	"api-gateway/clients"
	"net/http"

	"github.com/recktt77/projectProto-definitions/gen/project_service/genproto/project"

	"github.com/gin-gonic/gin"
)

func GetAllProjects(c *gin.Context) {
	clientID := c.Query("client_id")

	res, err := clients.GetProjectClient().GetAllProjects(c, &project.GetAllProjectsRequest{
		ClientId: clientID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func GetProject(c *gin.Context) {
	id := c.Param("id")

	res, err := clients.GetProjectClient().GetProject(c, &project.GetProjectRequest{
		Id: id,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func CreateProject(c *gin.Context) {
	var req project.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := clients.GetProjectClient().CreateProject(c, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func UpdateProject(c *gin.Context) {
	var req project.UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.ProjectId = c.Param("id")

	res, err := clients.GetProjectClient().UpdateProject(c, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func DeleteProject(c *gin.Context) {
	id := c.Param("id")

	res, err := clients.GetProjectClient().DeleteProject(c, &project.DeleteProjectRequest{
		Id: id,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func CreateTask(c *gin.Context) {
	projectID := c.Param("id")
	var req project.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.ProjectId = projectID

	res, err := clients.GetProjectClient().CreateTask(c, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func UpdateTask(c *gin.Context) {
	projectID := c.Param("id")
	var req project.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.ProjectId = projectID

	res, err := clients.GetProjectClient().UpdateTask(c, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func DeleteTask(c *gin.Context) {
	projectID := c.Param("id")
	var req project.DeleteTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.ProjectId = projectID

	res, err := clients.GetProjectClient().DeleteTask(c, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func GetTasks(c *gin.Context) {
	projectID := c.Param("id")

	res, err := clients.GetProjectClient().GetTasks(c, &project.GetTasksRequest{
		ProjectId: projectID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func AttachFile(c *gin.Context) {
	projectID := c.Param("id")
	var req project.AttachFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.ProjectId = projectID

	res, err := clients.GetProjectClient().AttachFile(c, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func DeleteFile(c *gin.Context) {
	projectID := c.Param("id")
	var req project.DeleteFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.ProjectId = projectID

	res, err := clients.GetProjectClient().DeleteFile(c, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
