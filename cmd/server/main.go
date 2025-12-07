package main

import (
	"money-buddy-backend/internal/handlers"
	"money-buddy-backend/internal/repositories"
	"money-buddy-backend/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	repo := repositories.NewExpenseRepositoryMemory()
	service := services.NewExpenseService(repo)
	handlers.NewExpenseHandler(r, service)

	r.Run() // デフォルトで:8080で起動
}
