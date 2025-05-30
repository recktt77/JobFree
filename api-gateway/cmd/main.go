package main

import (
	"log"
	"net/http"

	"api-gateway/router"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	router.SetupRoutes(r)

	log.Println("API Gateway listening on :8080")
	http.ListenAndServe(":8080", r)
}
