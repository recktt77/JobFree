package router

import (
	"api-gateway/internal/handlers"
	"api-gateway/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/")

	// Project endpoints
	api.GET("/projects", handlers.GetAllProjects)
	api.GET("/projects/:id", handlers.GetProject)
	api.POST("/projects", middleware.ValidateJWT, handlers.CreateProject)
	api.PUT("/projects/:id", middleware.ValidateJWT, handlers.UpdateProject)
	api.DELETE("/projects/:id", middleware.ValidateJWT, handlers.DeleteProject)

	api.GET("/projects/:id/tasks", handlers.GetTasks)
	api.POST("/projects/:id/tasks", middleware.ValidateJWT, handlers.CreateTask)
	api.PUT("/projects/:id/tasks", middleware.ValidateJWT, handlers.UpdateTask)
	api.DELETE("/projects/:id/tasks", middleware.ValidateJWT, handlers.DeleteTask)

	api.POST("/projects/:id/files", middleware.ValidateJWT, handlers.AttachFile)
	api.DELETE("/projects/:id/files", middleware.ValidateJWT, handlers.DeleteFile)

	// Review endpoints
	api.GET("/reviews", handlers.GetReviews)
	api.POST("/reviews", middleware.ValidateJWT, handlers.LeaveReview)
	api.POST("/reviews/moderate", middleware.ValidateJWT, handlers.ModerateReview)

	// Messaging endpoints
	api.GET("/messages/:conversation_id", middleware.ValidateJWT, handlers.GetMessages)
	api.POST("/messages", middleware.ValidateJWT, handlers.SendMessage)

	// Auth endpoints
	api.POST("/auth/register", handlers.RegisterUser)
	api.POST("/auth/login", handlers.LoginUser)
	api.POST("/auth/profile", middleware.ValidateJWT, handlers.GetProfile)
	api.PUT("/auth/profile", middleware.ValidateJWT, handlers.UpdateProfile)

	api.POST("/payments", middleware.ValidateJWT, handlers.CreatePayment)
	api.POST("/payments/get", middleware.ValidateJWT, handlers.GetPayment)
	api.POST("/payments/list", middleware.ValidateJWT, handlers.ListUserPayments)

	api.POST("/subscriptions", middleware.ValidateJWT, handlers.Subscribe)
	api.POST("/subscriptions/cancel", middleware.ValidateJWT, handlers.CancelSubscription)
	api.POST("/subscriptions/status", middleware.ValidateJWT, handlers.GetSubscriptionStatus)
	api.POST("/subscriptions/all", middleware.ValidateJWT, handlers.GetSubscriptions)

	// Plan endpoints
	api.POST("/plans", middleware.ValidateJWT, handlers.CreatePlan)
	api.POST("/plans/get", middleware.ValidateJWT, handlers.GetPlanByID)
	api.PUT("/plans", middleware.ValidateJWT, handlers.UpdatePlan)
	api.POST("/plans/list", middleware.ValidateJWT, handlers.GetListOfPlans)

	// bids
	// Matching endpoints
	api.POST("/bids", middleware.ValidateJWT, handlers.CreateBid)
	api.GET("/bids/:project_id", handlers.GetBidsForProject)
	api.POST("/match", middleware.ValidateJWT, handlers.MatchFreelancers)

	// Admin
	api.POST("/admin/ban", middleware.ValidateJWT, handlers.BanUser)
	api.POST("/admin/delete-review", middleware.ValidateJWT, handlers.DeleteReview)
	api.POST("/admin/moderate-project", middleware.ValidateJWT, handlers.ModerateProject)
	api.GET("/admin/stats", middleware.ValidateJWT, handlers.GetPlatformStats)
}
