package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Muhammad-Magomedov/blog/internal/model"
	"github.com/Muhammad-Magomedov/blog/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type Handler struct {
	service service.BlogService
}

func New(service service.BlogService) Handler {
	return Handler{
		service: service,
	}
}

func (h *Handler) CreateUser(c *gin.Context) {
	var req model.CreateUserReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, "У вас невалидный запрос")
		return
	}

	err = h.service.CreateUser(c, req)
	if err != nil {
		log.Println("create user handler: %w", err)
		c.JSON(http.StatusInternalServerError, "Попробуйте позже")
	}

	c.Status(http.StatusOK)
}

func (h *Handler) CreatePost(c *gin.Context) {
	var req model.CreatePostReq
	id := c.Param("id")
	userId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "У вас невалидный запрос")
		return
	}

	err = h.service.CreatePost(c, userId, req)
	if err != nil {
		log.Println("create post handler: %w", err)
		c.JSON(http.StatusInternalServerError, "Попробуйте позже")
	}

	c.Status(http.StatusOK)
}

func (h *Handler) CreateComment(c *gin.Context) {
	var req model.CreateCommentReq
	stringUserId := c.Param("userId")
	stringId := c.Param("id")
	userId, err := strconv.Atoi(stringUserId)
	if err != nil {
		c.JSON(http.StatusBadRequest, "У вас невалидный запрос")
		return
	}
	id, err := strconv.Atoi(stringId)
	if err != nil {
		c.JSON(http.StatusBadRequest, "У вас невалидный запрос")
		return
	}

	fmt.Println("TEEEEEST", req.Body)
	err = h.service.CreateComment(c, userId, id, req)
	if err != nil {
		log.Println("create comment handler: %w", err)
		c.JSON(http.StatusInternalServerError, "Попробуйте позже")
	}

	c.Status(http.StatusOK)
}

func (h *Handler) GetUser(c *gin.Context) {
	stringId := c.Param("id")
	id, err := strconv.Atoi(stringId)
	if err != nil {
		log.Println("Error in GetUser handler: %w", err)
		c.JSON(http.StatusBadRequest, "Невалидные данные, пользователь не найден")
		return
	}

	user, err := h.service.GetUser(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Такой пользователь не найден")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"name":     user.Name,
		"is_admin": user.IsAdmin,
		"email":    user.Email,
	})
}

func (h *Handler) GetPost(c *gin.Context) {
	stringId := c.Param("id")
	id, err := strconv.Atoi(stringId)
	if err != nil {
		log.Println("Error in GetPost handler: %w", err)
		c.JSON(http.StatusBadRequest, "Невалидные данные, пост не найден")
		return
	}

	post, err := h.service.GetPost(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Такой пост не найден")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"title": post.Title,
		"body":  post.Body,
		"views": post.Views,
	})
}

func (h *Handler) GetPosts(c *gin.Context) {
	posts, err := h.service.GetPosts(c)
	if err != nil {
		log.Println("get posts handler: %w", err)
		c.JSON(http.StatusInternalServerError, "Не удалось получить посты")
		return
	}

	c.JSON(http.StatusOK, posts)
}

func (h *Handler) GetCommentsByPostId(c *gin.Context) {
	id := c.Param("id")
	postId, err := strconv.Atoi(id)
	if err != nil {
		log.Println("get comments handler: %w", err)
		c.JSON(http.StatusInternalServerError, "Не удалось получить комментарии")
		return
	}
	comments, err := h.service.GetCommentsByPostId(c, postId)
	if err != nil {
		log.Println("error service.GetCommentsByPostId: %w", err)
		c.JSON(http.StatusInternalServerError, "Комментарии не найдены")
		return
	}

	c.JSON(http.StatusOK, comments)
}

func (h *Handler) GetUsers(c *gin.Context) {
	users, err := h.service.GetUsers(c)
	if err != nil {
		log.Println("get users handler: %w", err)
		c.JSON(http.StatusInternalServerError, "Не удалось получить пользователей")
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *Handler) GetPostsByUserID(c *gin.Context) {
	param := c.Param("id")
	userId, err := strconv.Atoi(param)
	if err != nil {
		log.Println("Error in hanlder GetPostsByUserId problem with id: %w", err)
		c.JSON(http.StatusBadRequest, "Некорректный айди")
	}

	posts, err := h.service.GetPostsByUserID(c, userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Println("Error in GetPostsByUserId with no rows: %w", err)
			c.JSON(http.StatusBadRequest, "Посты не найдены")
			return
		}
		log.Println("Error in GetPostsByUserId: %w", err)
		c.JSON(http.StatusBadGateway, "Произошла какая-то ошибка, повторите позже")
		return
	}

	c.JSON(http.StatusOK, posts)
}
