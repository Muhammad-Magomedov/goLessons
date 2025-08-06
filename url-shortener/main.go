package main

import (
	"context"
	"log"
	"url-shortener/handler"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func main() {
	connString := "postgres://admin:admin@localhost:5432/links"
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		log.Fatal("Ошибка при подключении к БД: ", err)
	}
	r := gin.Default()

	linksHandler := handler.NewHandler(conn)

	r.POST("/shorten", linksHandler.CreateLink)
	r.GET("/:path", linksHandler.Redirect)
	// r.GET("/posts/:id", postsHandler.GetPostById)
	// r.DELETE("/posts/:id", postsHandler.DeletePost)
	// r.PATCH("/posts/:id", postsHandler.UpdatePost)

	r.Run(":8085")
}
