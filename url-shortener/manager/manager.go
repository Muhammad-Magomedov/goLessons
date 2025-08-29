package manager

import (
	"context"
	"fmt"
	"url-shortener/cache"
	"url-shortener/repo"
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

func (h *Manager) CachePopularLinks(linksRepository *repo.Repository, linksCache *cache.LinksCache) error {
	links, err := linksRepository.GetPopularLinks(context.Background(), 10)
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
