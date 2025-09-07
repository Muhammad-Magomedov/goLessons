package main

import (
	"context"
	"log"
	"time"

	"github.com/Muhammad-Magomedov/blog/internal/handler"
	"github.com/Muhammad-Magomedov/blog/internal/repo"
	"github.com/Muhammad-Magomedov/blog/internal/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func main() {
	connString := "postgres://postgres:postgres@localhost:5432/postgres"
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		log.Fatal("Ошибка при подключении к БД: ", err)
	}

	blogRepository := repo.New(conn)
	blogService := service.New(*blogRepository)
	blogHandler := handler.New(*blogService)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/user/:id", blogHandler.GetUser)
	r.POST("/user", blogHandler.CreateUser)
	r.GET("/users", blogHandler.GetUsers)
	r.Run(":8080")
}
