package manager

import (
	"context"
	"fmt"
	"url-shortener/cache"
	"url-shortener/repo"

	"github.com/gin-gonic/gin"
)

type Manager struct {
	LinksRepository *repo.Repository
	LinksCache      *cache.LinksCache
}

func ManagerHandler(linksRepo *repo.Repository, linksCache *cache.LinksCache) Manager {
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

func (h *Manager) CachePopularLinks(linksRepository *repo.Repository, linksCache *cache.LinksCache) error {
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
