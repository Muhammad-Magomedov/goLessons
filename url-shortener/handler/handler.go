package handler

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

const HostURL = "127.0.0.1:8085/"

type Handler struct {
	db *pgx.Conn
}

type CreateLinkRequest struct {
	Link string `json:"link"`
}

type LinkResponse struct {
	LongLink  string `json:"long_link"`
	ShortLink string `json:"short_link"`
}

func NewHandler(db *pgx.Conn) Handler {
	return Handler{
		db: db,
	}
}

func (h *Handler) CreateLink(c *gin.Context) {
	var req CreateLinkRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, "У вас невалидный запрос")
		return
	}

	b := make([]byte, 6)
	rand.Read(b)
	shortLink := base64.URLEncoding.EncodeToString(b)[:6]
	longLink := LinkResponse{}
	err = h.db.QueryRow(c, "SELECT long_link, short_link FROM links WHERE long_link=$1", req.Link).Scan(&longLink.LongLink, &longLink.ShortLink)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"short": longLink.ShortLink,
			"long":  longLink.LongLink,
		})
		return
	} else if err != pgx.ErrNoRows {
		c.JSON(http.StatusInternalServerError, "Ошибка базы данных")
		return
	}

	err = h.db.QueryRow(c, "select short_link from links where short_link=$1", shortLink).Scan()
	if err != pgx.ErrNoRows {
		c.JSON(http.StatusInternalServerError, "Ошибка базы данных")
		return
	}

	_, err = h.db.Exec(c, "insert into links (long_link, short_link) VALUES ($1, $2)", req.Link, shortLink)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Произошла ошибка, попробуйте позже")
	}

	c.JSON(http.StatusOK, gin.H{
		"short": HostURL + shortLink,
		"long":  req.Link,
	})
}

func (h *Handler) Redirect(c *gin.Context) {
	shortLink := c.Param("path")
	var longLink string

	row := h.db.QueryRow(c, "select long_link from links where short_link=$1", shortLink)
	err := row.Scan(&longLink)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, "Ссылка не найдена")
			return
		}
		c.JSON(http.StatusInternalServerError, "Произошла ошибка, попробуйте позже")
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, longLink)
}
