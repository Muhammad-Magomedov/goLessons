package handler

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"
	"net/http"
	"time"
	"unicode"
	"url-shortener/cache"
	"url-shortener/manager"
	"url-shortener/repo"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

const HostURL = "127.0.0.1:8080/"
const MaxLinkLength = 6
const MinLinkLength = 3

type Handler struct {
	LinksRepository *repo.Repository
	LinksCache      *cache.LinksCache
	LinksManager    *manager.Manager
}

type CreateLinkRequest struct {
	Link       string  `json:"link"`
	CustomLink *string `json:"custom_link"`
}

type LinkResponse struct {
	LongLink  string `json:"long_link"`
	ShortLink string `json:"short_link"`
}

func NewHandler(linksRepo *repo.Repository, linksCache *cache.LinksCache) Handler {
	return Handler{
		LinksRepository: linksRepo,
		LinksCache:      linksCache,
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

	if req.CustomLink != nil {
		customLink := *req.CustomLink

		err := validateShortLink(customLink)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Невалидная кастомная ссылка. Длина должна быть 6 символов. Символы могут быть цифрами и английскими буквами"})
			return
		}

		isExists, err := h.LinksRepository.IsShortExists(c, customLink)
		if err != nil {
			log.Printf("error IsShortExists: %w", err)
			c.JSON(http.StatusInternalServerError, "Произошла ошибка БД, попробуйте позже")
			return
		}

		if isExists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Такая ссылка занята. Попробуйте другую"})
			return
		}

		err = h.LinksRepository.CreateLink(c, req.Link, customLink)
		if err != nil {
			log.Printf("error CreateLink: %w", err)
			c.JSON(http.StatusInternalServerError, "Произошла ошибка, попробуйте позже")
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"short": HostURL + customLink,
			"long":  req.Link,
		})
		return

	}

	shortLink := ""

	for {
		b := make([]byte, MaxLinkLength)
		rand.Read(b)
		shortLink = base64.URLEncoding.EncodeToString(b)[:MaxLinkLength]

		isExists, err := h.LinksRepository.IsShortExists(c, shortLink)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "Произошла ошибка БД, попробуйте позже")
			return
		}

		if !isExists {
			break
		}
	}

	err = h.LinksRepository.CreateLink(c, req.Link, shortLink)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Произошла ошибка, попробуйте позже")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"short": HostURL + shortLink,
		"long":  req.Link,
	})
}

func (h *Handler) Redirect(c *gin.Context) {
	shortLink := c.Param("path")

	start := time.Now()
	longLink, err := h.LinksRepository.GetLongByShort(c, shortLink)
	if err != nil {
		log.Printf("error LinksCache.GetLink: ", err)
	}

	if longLink == "" {
		longLink, err = h.LinksRepository.GetLongByShort(c, shortLink)
		log.Println("vid db: ", time.Since(start))
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				c.JSON(http.StatusNotFound, "Ссылка не найдена")
				return
			}
			log.Println("error GetLongByShort: ", err)
			c.JSON(http.StatusInternalServerError, "Произошла ошибка, попробуйте позже")
		}
	} else {
		log.Println("via cache: ", time.Since(start))
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

func validateShortLink(link string) error {
	if len(link) < MinLinkLength {
		return errors.New("error, link too short")
	}

	for _, r := range link {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return errors.New("error invalid symbol in link")
		}
	}

	return nil
}
