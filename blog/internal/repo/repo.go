package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type Post struct {
	Id        int       `json:"id"`
	UserId    int       `json:"user_id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Views     int       `json:"views"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Comments struct {
	Id        int       `json:"id"`
	UserId    int       `json:"user_id"`
	PostId    int       `json:"post_id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Repository struct {
	db *pgx.Conn
}

func New(db *pgx.Conn) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreatePost(ctx context.Context, user_id int, title string, body string, views int) (string, error) {
	_, err := r.db.Exec(ctx, "insert into posts (user_id, title, body, views) VALUES ($1, $2, $3, $4)", user_id, title, body, views)
	if err != nil {
		return "", fmt.Errorf("error in CreatePost: %w", err)
	}

	return "Пост успешно создан", nil
}

func (r *Repository) GetPost(ctx context.Context, id int) (Post, error) {
	var post Post
	err := r.db.QueryRow(ctx, "select id, user_id, title, body, views, created_at, updated_at from posts where id=$1", id).Scan(&post.Id, &post.UserId, &post.Title, &post.Body, &post.Views, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return Post{}, fmt.Errorf("error in GetPost: %w", err)
	}

	return post, nil
}

func (r *Repository) RemovePost(ctx context.Context, id int) (string, error) {
	_, err := r.db.Exec(ctx, "delete from posts where id=$1", id)
	if err != nil {
		return "", fmt.Errorf("error in RemovePost: %w", err)
	}

	return "Пост удален", nil
}

func (r *Repository) UpdatePost(ctx context.Context, id int, title string, body string) (string, error) {
	_, err := r.db.Exec(ctx, "update posts set title=$1, body=$2 where id=$3", title, body, id)
	if err != nil {
		return "", fmt.Errorf("error in UpdatePost: %w", err)
	}

	return "Пост обновлен", nil
}

func (r *Repository) CreateUser(ctx context.Context, name string, email string, is_admin bool) (string, error) {
	_, err := r.db.Exec(ctx, "insert into users (name, email, is_admin) VALUES ($1, $2, $3)", name, email, is_admin)
	if err != nil {
		return "", fmt.Errorf("error in CreateUser: %w", err)
	}

	return "Пост успешно создан", nil
}

func (r *Repository) RemoveUser(ctx context.Context, id int) (string, error) {
	_, err := r.db.Exec(ctx, "delete from users where id=$1", id)
	if err != nil {
		return "", fmt.Errorf("error in RemoveUser: %w", err)
	}

	return "Пост удален", nil
}

func (r *Repository) UpdateUser(ctx context.Context, id int, name string, email string, is_admin bool) (string, error) {
	_, err := r.db.Exec(ctx, "update posts set name=$1, email=$2, is_admin=$3 where id=$4", name, email, is_admin, id)
	if err != nil {
		return "", fmt.Errorf("error in UpdateUser: %w", err)
	}

	return "Пост обновлен", nil
}

func (r *Repository) CreateComment(ctx context.Context, user_id int, post_id int, body string) (string, error) {
	_, err := r.db.Exec(ctx, "insert into comments (user_id, post_id, body) VALUES ($1, $2, $3)", user_id, post_id, body)
	if err != nil {
		return "", fmt.Errorf("error in CreateComment: %w", err)
	}
	return "Комментарий успешно создан", nil
}

func (r *Repository) RemoveComment(ctx context.Context, id int) (string, error) {
	_, err := r.db.Exec(ctx, "delete from comments where id=$1", id)
	if err != nil {
		return "", fmt.Errorf("error in RemoveComment: %w", err)
	}
	return "Комментарий удалён", nil
}

func (r *Repository) UpdateComment(ctx context.Context, id int, body string) (string, error) {
	_, err := r.db.Exec(ctx, "update comments set body=$1 where id=$2", body, id)
	if err != nil {
		return "", fmt.Errorf("error in UpdateComment: %w", err)
	}
	return "Комментарий обновлён", nil
}

func (r *Repository) GetComment(ctx context.Context, id int) (Comments, error) {
	var comment Comments
	err := r.db.QueryRow(ctx, "select id, user_id, post_id, body, created_at, updated_at from comments where id=$1", id).
		Scan(&comment.Id, &comment.UserId, &comment.PostId, &comment.Body, &comment.CreatedAt, &comment.UpdatedAt)
	if err != nil {
		return Comments{}, fmt.Errorf("error in GetComment: %w", err)
	}
	return comment, nil
}
