package main

import (
	"money-buddy-backend/internal/db"
	"money-buddy-backend/internal/handlers"
	"money-buddy-backend/internal/repositories"
	"money-buddy-backend/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	dbConn, err := db.NewDB()
	if err != nil {
		panic(err)
	}

	repo := repositories.NewExpenseRepositoryPostgres(dbConn)
	service := services.NewExpenseService(repo)
	handlers.NewExpenseHandler(r, service)

	r.Run() // デフォルトで:8080で起動
}
