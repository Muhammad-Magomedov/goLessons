package main

import (
	"crud/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	postsHandler := handler.Handler{
		LastID: 0,
		Posts:  make(map[int]handler.Post),
	}

	r.POST("/posts", postsHandler.CreatePostHandler)
	r.PATCH("/posts/:id", postsHandler.UpdatePostHandler)
	r.PUT("/posts/:id", postsHandler.UpdatePostHandler)
	r.GET("/posts", postsHandler.GetPostsHanler)
	r.GET("/posts/:id", postsHandler.GetPostByIdHandler)
	r.DELETE("/posts/:id", postsHandler.DeletePostHandler)

	r.Run(":8085")
}
