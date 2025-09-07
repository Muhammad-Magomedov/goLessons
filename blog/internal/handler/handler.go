package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Muhammad-Magomedov/blog/internal/model"
	"github.com/Muhammad-Magomedov/blog/internal/service"
	"github.com/gin-gonic/gin"
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

func (h *Handler) GetUsers(c *gin.Context) {
	users, err := h.service.GetUsers(c)
	if err != nil {
		log.Println("get users handler: %w", err)
		c.JSON(http.StatusInternalServerError, "Не удалось получить пользователей")
		return
	}

	c.JSON(http.StatusOK, users)
}
