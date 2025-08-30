package manager

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
	"url-shortener/cache"
	"url-shortener/repo"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx"
)

type Manager struct {
	LinksRepository *repo.Repository
	LinksCache      *cache.LinksCache
}

func New(linksRepo *repo.Repository, linksCache *cache.LinksCache) Manager {
	return Manager{
		LinksRepository: linksRepo,
		LinksCache:      linksCache,
	}
}

func (m *Manager) CreateLink(c *gin.Context, longLink string, shortLink string) error {
	return m.LinksRepository.CreateLink(c, longLink, shortLink)
}

func (m *Manager) GetLongByShort(c *gin.Context, shortLink string) (string, error) {
	return m.LinksRepository.GetLongByShort(c, shortLink)
}

func (m *Manager) GetShortByLong(c *gin.Context, longLink string) (string, error) {
	return m.LinksRepository.GetShortByLong(c, longLink)
}

func (m *Manager) RemoveLink(c *gin.Context, shortLink string) (string, error) {
	return m.LinksRepository.RemoveLink(c, shortLink)
}

func (m *Manager) IsShortExists(c *gin.Context, shortLink string) (bool, error) {
	return m.LinksRepository.IsShortExists(c, shortLink)
}

func (m *Manager) Redirect(c *gin.Context, shortLink string) (string, error) {
	return m.LinksRepository.Redirect(c, shortLink)
}

func (m *Manager) StoreRedirect(c *gin.Context, params repo.StoreRedirectParams) error {
	return m.LinksRepository.StoreRedirect(c, params)
}

func (m *Manager) GetRedirectsByShortLink(c *gin.Context, shortLink string) ([]repo.Redirect, error) {
	return m.LinksRepository.GetRedirectsByShortLink(c, shortLink)
}

func (m *Manager) GetPopularLinks(ctx context.Context, n int) ([]repo.LinkPair, error) {
	return m.LinksRepository.GetPopularLinks(ctx, n)
}

func (m *Manager) CachePopularLinks(linksRepository *repo.Repository, linksCache *cache.LinksCache) error {
	links, err := linksRepository.GetPopularLinks(context.Background(), 2)
	if err != nil {
		return fmt.Errorf("error updateCache GetPopularLinks: %w", err)
	}

	for _, link := range links {
		err := linksCache.StoreLink(link.Short, link.Long)
		if err != nil {
			return fmt.Errorf("error updateCache StoreLink: %w", err)
		}
	}

	return nil
}

func (m *Manager) FindLink(shortLink string, c *gin.Context) (string, error) {
	start := time.Now()
	longLink, err := m.LinksCache.GetLink(shortLink)
	if err != nil {
		log.Printf("error FindLink: ", err)
	}

	log.Println(longLink)

	if longLink == "" {
		longLink, err = m.LinksRepository.GetLongByShort(c, shortLink)
		log.Println("vid db: ", time.Since(start))
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				c.JSON(http.StatusNotFound, "Ссылка не найдена")
				return "", err
			}
			log.Println("error GetLongByShort: ", err)
			c.JSON(http.StatusInternalServerError, "Произошла ошибка, попробуйте позже")
		}
	} else {
		log.Println("via cache: ", time.Since(start))
	}

	return longLink, nil
}
