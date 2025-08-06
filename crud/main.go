package main

import (
	"context"
	"crud/handler"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func main() {
	connString := "postgres://admin:admin@localhost:5432/mydb"
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		log.Fatal("Ошибка при подключении к БД: ", err)
	}
	r := gin.Default()

	postsHandler := handler.NewHandler(conn)

	r.POST("/posts", postsHandler.CreatePost)
	r.GET("/posts", postsHandler.GetPosts)
	r.GET("/posts/:id", postsHandler.GetPostById)
	r.DELETE("/posts/:id", postsHandler.DeletePost)
	r.PATCH("/posts/:id", postsHandler.UpdatePost)

	r.Run(":8085")
}
