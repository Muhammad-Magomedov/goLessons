package service

import (
	"context"
	"fmt"

	"github.com/Muhammad-Magomedov/blog/internal/model"
	"github.com/Muhammad-Magomedov/blog/internal/repo"
	"golang.org/x/crypto/bcrypt"
)

type BlogService struct {
	blogManager repo.Repository
}

func New(blogManager repo.Repository) *BlogService {
	return &BlogService{
		blogManager: blogManager,
	}
}

func (b *BlogService) CreateUser(ctx context.Context, user model.CreateUserReq) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error bcrypt.GenerateFromPassword: %w", err)
	}

	err = b.blogManager.CreateUser(ctx, repo.CreateUser{
		Name:           user.Name,
		HashedPassword: string(hashedPassword),
		Email:          user.Email,
	})
	if err != nil {
		return fmt.Errorf("error repo.CreateUser: %w", err)
	}

	return nil
}

func (b *BlogService) GetUser(ctx context.Context, id int) (model.User, error) {
	user, err := b.blogManager.GetUser(ctx, id)
	if err != nil {
		return model.User{}, fmt.Errorf("error repo.GetUser: %w", err)
	}

	return user, nil
}

func (b *BlogService) GetUsers(ctx context.Context) ([]model.User, error) {
	users, err := b.blogManager.GetUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("service.GetUsers: %w", err)
	}
	return users, nil
}
