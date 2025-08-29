package main

import (
	"context"
	"log"
	"time"

	"url-shortener/cache"
	"url-shortener/handler"
	"url-shortener/manager"
	"url-shortener/repo"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/jackc/pgx/v5"
)

const cacheLinksInterval = time.Hour

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := rdb.Ping().Result()
	if err != nil {
		log.Fatalf("Ошибка подключения к Redis: %v", err)
	}

	connString := "postgres://postgres:postgres@localhost:5432/shortener-db"
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		log.Fatal("Ошибка при подключении к БД: ", err)
	}

	linksRepository := repo.NewRepository(conn)
	linksCache := cache.New(rdb)
	linksHandler := handler.NewHandler(linksRepository, linksCache)
	linksManager := manager.ManagerHandler(linksRepository, linksCache)

	go func() {
		err := linksManager.CachePopularLinks(linksRepository, linksCache)
		if err != nil {
			log.Println(err)
		}

		c := time.Tick(cacheLinksInterval)
		for range c {
			err := linksManager.CachePopularLinks(linksRepository, linksCache)
			if err != nil {
				log.Println(err)
			}
		}
	}()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	r.POST("/shorten", linksHandler.CreateLink)
	r.GET("/:path", linksHandler.Redirect)
	r.GET("/analytics/:path", linksHandler.GetAnalytics)
	r.Run(":8080")
}
