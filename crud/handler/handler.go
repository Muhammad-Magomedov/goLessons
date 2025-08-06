package handler

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type Post struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

type UpdatePostRequest struct {
	Title *string `json:"title"`
	Body  *string `json:"body"`
}

type Handler struct {
	db *pgx.Conn
}

func NewHandler(db *pgx.Conn) Handler {
	return Handler{
		db: db,
	}
}

func checkId(idStr string, c *gin.Context) int {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return -1
	}
	return id
}

func (h *Handler) CreatePost(c *gin.Context) {
	var post Post
	err := c.BindJSON(&post)
	if err != nil {
		c.JSON(http.StatusBadRequest, "У вас невалидный запрос")
		return
	}

	_, err = h.db.Exec(c, "insert into posts (title, body) values ($1, $2)", post.Title, post.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Что-то пошло не так")
		log.Println("error insert into posts: ", err)
	}

	c.JSON(http.StatusOK, post)
}

func (h *Handler) GetPosts(c *gin.Context) {
	posts, err := h.db.Query(c, "select title, body, id from posts")
	if err != nil {
		log.Fatalf("Ошибка запроса: %v", err)
	}

	defer posts.Close()

	sum := []Post{}

	for posts.Next() {
		var post Post
		err = posts.Scan(&post.Title, &post.Body, &post.ID)
		if err != nil {
			log.Println("bbbbbbb", err)
			return
		}
		sum = append(sum, post)
	}

	if posts.Err() != nil {
		log.Println("aaaaa")
		return
	}

	c.JSON(http.StatusOK, sum)
}

func (h *Handler) GetPostById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Невалидный id")
		return
	}

	post := Post{}
	err = h.db.QueryRow(c, "select * from posts where id=$1", id).Scan(&post.CreatedAt, &post.ID, &post.Body, &post.Title)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Пост не найден")
		return
	}
	c.JSON(http.StatusOK, post)
}

func (h *Handler) DeletePost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Невалидный id")
		return
	}

	_, err = h.db.Exec(c, "delete from posts where id=$1", id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Не удалось удалить пост")
		return
	}

	c.JSON(http.StatusOK, "Запись удалена")
}

func (h *Handler) UpdatePost(c *gin.Context) {
	idStr := c.Param("id")
	id := checkId(idStr, c)

	post := Post{}
	err := h.db.QueryRow(c, "select * from posts where id=$1", id).Scan(&post.CreatedAt, &post.ID, &post.Body, &post.Title)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Пост не найден")
		return
	}

	var updatePostRequest UpdatePostRequest
	err = c.BindJSON(&updatePostRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, "У вас невалидный запрос")
		return
	}

	_, err = h.db.Exec(c, "update posts set title=$1, body=$2 where id=$3", updatePostRequest.Title, updatePostRequest.Body, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Не удалось изменить пост")
		return
	}

	c.JSON(http.StatusOK, "Запись изменена")
}
