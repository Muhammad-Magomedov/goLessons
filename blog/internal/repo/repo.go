package repo

import (
	"context"
	"fmt"

	"github.com/Muhammad-Magomedov/blog/internal/model"
	"github.com/jackc/pgx/v5"
)

type CreateComment struct {
	UserId int    `json:"user_id"`
	PostId int    `json:"post_id"`
	Body   string `json:"body"`
}

type CreatePost struct {
	UserId int    `json:"user_id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

type CreateUser struct {
	Name           string
	HashedPassword string
	Email          string
	IsAdmin        bool
}

type UpdateComment struct {
	//Здесь нужен будет пост айди?
	Id   int    `json:"id"`
	Body string `json:"body"`
}

type Repository struct {
	db *pgx.Conn
}

func New(db *pgx.Conn) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreatePost(ctx context.Context, post CreatePost) error {
	_, err := r.db.Exec(ctx, `
	insert into posts
	 (user_id, title, body)
	  VALUES ($1, $2, $3)`,
		&post.UserId, &post.Title, &post.Body)
	if err != nil {
		return fmt.Errorf("error in CreatePost: %w", err)
	}

	return nil
}

func (r *Repository) GetPost(ctx context.Context, id int) (model.Post, error) {
	var post model.Post
	err := r.db.QueryRow(
		ctx,
		`select
				id,
				user_id,
				title,
				body,
				views,
				created_at,
				updated_at
		   from posts
		  where id=$1`,
		id,
	).Scan(
		&post.Id,
		&post.UserId,
		&post.Title,
		&post.Body,
		&post.Views,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		return model.Post{}, fmt.Errorf("error in GetPost: %w", err)
	}

	return post, nil
}

func (r *Repository) GetPosts(ctx context.Context) ([]model.Post, error) {
	var posts []model.Post
	rows, err := r.db.Query(
		ctx,
		`select
				id,
				user_id,
				title,
				body,
				views,
				created_at,
				updated_at
		   from posts`,
	)
	if err != nil {
		return nil, fmt.Errorf("repo.GetPosts: %w", err)
	}

	for rows.Next() {
		var post model.Post
		err := rows.Scan(&post.Id,
			&post.UserId,
			&post.Title,
			&post.Body,
			&post.Views,
			&post.CreatedAt,
			&post.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("repo.GetPosts scan: %w", err)
		}
		posts = append(posts, post)
	}

	return posts, nil
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

func (r *Repository) GetUser(ctx context.Context, id int) (model.User, error) {
	var user model.User
	err := r.db.QueryRow(
		ctx,
		`select
				id,
				name,
				email,
				is_admin,
				created_at,
				updated_at
		   from users
		  where id=$1`,
		id,
	).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.IsAdmin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return model.User{}, fmt.Errorf("error in GetUser: %w", err)
	}
	return user, nil
}

func (r *Repository) GetUsers(ctx context.Context) ([]model.User, error) {
	rows, err := r.db.Query(ctx, "SELECT id, name, email, is_admin FROM users")
	if err != nil {
		return nil, fmt.Errorf("repo.GetUsers: %w", err)
	}

	var users []model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.IsAdmin)
		if err != nil {
			return nil, fmt.Errorf("repo.GetUsers scan: %w", err)
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *Repository) CreateUser(ctx context.Context, user CreateUser) error {
	_, err := r.db.Exec(
		ctx,
		"insert into users (name, hashed_password, email, is_admin) VALUES ($1, $2, $3, $4)",
		user.Name,
		user.HashedPassword,
		user.Email,
		user.IsAdmin,
	)
	if err != nil {
		return fmt.Errorf("error in CreateUser: %w", err)
	}

	return nil
}

func (r *Repository) RemoveUser(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, "delete from users where id=$1", id)
	if err != nil {
		return fmt.Errorf("error in RemoveUser: %w", err)
	}

	return nil
}

type UpdateUser struct {
	Name    *string
	Email   *string
	IsAdmin *bool
}

func (r *Repository) UpdateUser(ctx context.Context, id int, user UpdateUser) (model.User, error) {
	var updatedUser model.User
	err := r.db.QueryRow(
		ctx,
		`update user
			set name=COALESCE($1, name),
				email=COALESCE($2, email),
				is_admin=COALESCE($3, is_admin),
				updated_at=now()
	      where id=$4
	  returning id,
				name,
				email,
				is_admin,
				created_at,
				updated_at;`,
		user.Name,
		user.Email,
		user.IsAdmin,
		id,
	).Scan(
		&updatedUser.Id,
		&updatedUser.Name,
		&updatedUser.Email,
		&updatedUser.IsAdmin,
		&updatedUser.CreatedAt,
		&updatedUser.UpdatedAt,
	)
	if err != nil {
		return model.User{}, fmt.Errorf("error in UpdateUser: %w", err)
	}

	return updatedUser, nil
}

func (r *Repository) CreateComment(ctx context.Context, comment CreateComment) error {
	_, err := r.db.Exec(ctx, "insert into comments (user_id, post_id, body) VALUES ($1, $2, $3)", comment.UserId, comment.PostId, comment.Body)
	if err != nil {
		return fmt.Errorf("error in CreateComment: %w", err)
	}
	return nil
}

func (r *Repository) RemoveComment(ctx context.Context, id int) (string, error) {
	_, err := r.db.Exec(ctx, "delete from comments where id=$1", id)
	if err != nil {
		return "", fmt.Errorf("error in RemoveComment: %w", err)
	}
	return "Комментарий удалён", nil
}

func (r *Repository) UpdateComment(ctx context.Context, id int, comment UpdateComment) (string, error) {
	var updatedComment UpdateComment
	err := r.db.QueryRow(ctx, "update comments set body=coalesce($1, body) where id=$2 returning id, body", comment.Body, id).Scan(&updatedComment.Body, &updatedComment.Id)
	if err != nil {
		return "", fmt.Errorf("error in UpdateComment: %w", err)
	}
	return "Комментарий обновлён", nil
}

func (r *Repository) GetCommentsByPostId(ctx context.Context, postId int) ([]model.Comment, error) {
	rows, err := r.db.Query(ctx, "select id, user_id, post_id, body, created_at, updated_at from comments where post_id=$1", postId)
	if err != nil {
		return nil, fmt.Errorf("repo.GetUsers: %w", err)
	}

	var comments []model.Comment
	for rows.Next() {
		var comment model.Comment
		err := rows.Scan(&comment.Id, &comment.UserId, &comment.PostId, &comment.Body, &comment.CreatedAt, &comment.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("repo.GetComments scan: %w", err)
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func (r *Repository) GetPostsByUserID(ctx context.Context, userId int) ([]model.Post, error) {
	rows, err := r.db.Query(ctx, "select id, user_id, title, body, views, created_at, updated_at from posts where user_id=$1", userId)
	if err != nil {
		return nil, fmt.Errorf("repo.GetUsers: %w", err)
	}

	var posts []model.Post
	for rows.Next() {
		var post model.Post
		err := rows.Scan(&post.Id, &post.UserId, &post.Title, &post.Body, &post.Views, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("repo.GetPostsByUserID scan: %w", err)
		}
		posts = append(posts, post)
	}
	return posts, nil
}
