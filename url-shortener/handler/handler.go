package handler

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"
	"net/http"
	"url-shortener/repo"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

const HostURL = "127.0.0.1:8080/"

type Handler struct {
	LinksRepository *repo.Repository
}

type CreateLinkRequest struct {
	Link       string  `json:"link"`
	CustomLink *string `json:"custom_link"`
}

type LinkResponse struct {
	LongLink  string `json:"long_link"`
	ShortLink string `json:"short_link"`
}

func NewHandler(linksRepo *repo.Repository) Handler {
	return Handler{
		LinksRepository: linksRepo,
	}
}

func (h *Handler) CreateLink(c *gin.Context) {
	var req CreateLinkRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, "У вас невалидный запрос")
		return
	}

	existingShortLink, err := h.LinksRepository.GetShortByLong(c, req.Link)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"short": HostURL + existingShortLink,
			"long":  req.Link,
		})
		return
	}
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		log.Println("Ошибка БД: ", err)
		c.JSON(http.StatusInternalServerError, "Ошибка базы данных")
		return
	}

	shortLink := ""
	for {
		b := make([]byte, 6)
		rand.Read(b)
		shortLink = base64.URLEncoding.EncodeToString(b)[:6]

		isExists, err := h.LinksRepository.IsShortExists(c, shortLink)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "Произошла ошибка БД, попробуйте позже")
			return
		}

		if !isExists {
			break
		}
	}

	var customLink string

	isExists, err := h.LinksRepository.IsCustomLinkExists(c, *req.CustomLink)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Произошла ошибка БД, попробуйте позже")
		return
	}

	if isExists {
		c.JSON(http.StatusBadRequest, "Такая ссылка уже существует")
		return
	}

	if req.CustomLink != nil {
		if len(*req.CustomLink) >= 3 {
			customLink = *req.CustomLink
		} else {
			c.JSON(http.StatusBadRequest, "Длина ссылки должна быть больше 2 символов")
			return
		}
	} else {
		customLink = shortLink
	}

	err = h.LinksRepository.CreateLink(c, req.Link, customLink)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Произошла ошибка, попробуйте позже")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"short": HostURL + customLink,
		"long":  req.Link,
	})
}

func (h *Handler) Redirect(c *gin.Context) {
	shortLink := c.Param("path")

	longLink, err := h.LinksRepository.GetLongByShort(c, shortLink)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, "Ссылка не найдена")
			return
		}
		c.JSON(http.StatusInternalServerError, "Произошла ошибка, попробуйте позже")
		return
	}

	err = h.LinksRepository.StoreRedirect(c, repo.StoreRedirectParams{
		UserAgent: c.GetHeader("User-Agent"),
		ShortLink: shortLink,
		LongLink:  longLink,
	})
	if err != nil {
		log.Printf("Ошибка при StoreRedirect: %v", err)
	}

	c.Redirect(http.StatusTemporaryRedirect, longLink)
}

func (h *Handler) GetAnalytics(c *gin.Context) {
	shortLink := c.Param("path")

	redirects, err := h.LinksRepository.GetRedirectsByShortLink(c, shortLink)
	if err != nil {
		log.Println("GetAnalytics", err)
		c.JSON(http.StatusInternalServerError, "Не удалось получить аналитику")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total_redirects": len(redirects),
		"redirects":       redirects,
	})
}
