package sites

import (
	"context"
	"errors"
	"fmt"
	"time"

	"dify-wp-sync/internal/logger"
	"dify-wp-sync/internal/redisstore"
)

const (
	sitesSetKey = "wp_sites" // A Redis set containing site IDs
)

type Manager struct {
	store *redisstore.RedisStore
}

func NewManager(store *redisstore.RedisStore) *Manager {
	return &Manager{store: store}
}

func (m *Manager) siteKey(siteID string) string {
	return fmt.Sprintf("wp_site:%s", siteID)
}

func (m *Manager) AddSite(ctx context.Context, cfg *SiteConfig) error {
	if cfg.PostDocMapping == nil {
		cfg.PostDocMapping = make(map[int]string)
	}
	err := m.store.SetJSON(ctx, m.siteKey(cfg.SiteID), cfg, 0)
	if err != nil {
		return err
	}
	return m.store.SAdd(ctx, sitesSetKey, cfg.SiteID)
}

func (m *Manager) GetSite(ctx context.Context, siteID string) (*SiteConfig, error) {
	var sc SiteConfig
	found, err := m.store.GetJSON(ctx, m.siteKey(siteID), &sc)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, errors.New("site not found")
	}
	return &sc, nil
}

func (m *Manager) UpdateSite(ctx context.Context, cfg *SiteConfig) error {
	return m.store.SetJSON(ctx, m.siteKey(cfg.SiteID), cfg, 0)
}

func (m *Manager) ListSites(ctx context.Context) ([]*SiteConfig, error) {
	ids, err := m.store.SMembers(ctx, sitesSetKey)
	if err != nil {
		return nil, err
	}
	var sites []*SiteConfig
	for _, id := range ids {
		sc, err := m.GetSite(ctx, id)
		if err != nil {
			logger.Log.Warnf("Error loading site %s: %v", id, err)
			continue
		}
		sites = append(sites, sc)
	}
	return sites, nil
}

func (m *Manager) UpdateLastSyncTime(ctx context.Context, siteID string, t time.Time) error {
	sc, err := m.GetSite(ctx, siteID)
	if err != nil {
		return err
	}
	if t.After(sc.LastSyncTime) {
		sc.LastSyncTime = t
		return m.UpdateSite(ctx, sc)
	}
	return nil
}

func (m *Manager) UpdatePostDocMapping(ctx context.Context, siteID string, postID int, docID string) error {
	sc, err := m.GetSite(ctx, siteID)
	if err != nil {
		return err
	}
	sc.PostDocMapping[postID] = docID
	return m.UpdateSite(ctx, sc)
}
