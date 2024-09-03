package sync

import (
	"fmt"
	"nexus-wallet/pkg/cache"
	"time"
)

type Params struct {
	SyncIntervalInMinutes int
}

type Manager struct {
	cacher cache.Cacher
	params Params
}

func NewManager(cacher cache.Cacher, params Params) *Manager {
	return &Manager{
		cacher: cacher,
		params: params,
	}
}

type LastSync struct {
	Timestamp time.Time `json:"timestamp"`
}

func (s Manager) NeedToSync(lastSync *LastSync) bool {
	return lastSync.Timestamp.IsZero() ||
		time.Now().After(lastSync.Timestamp.Add(time.Duration(s.params.SyncIntervalInMinutes)*time.Minute))
}

func (s Manager) GetLastSyncTime(baseKey string, mnemonicId int64) (*LastSync, error) {
	lastSync := LastSync{}
	err := s.cacher.Get(fmt.Sprintf("%s%d", baseKey, mnemonicId), &lastSync)
	if err != nil {
		return nil, fmt.Errorf("failed to get last sync: %s", err)
	}

	return &lastSync, nil
}

func (s Manager) SetLastSyncTime(baseKey string, mnemonicId int64) error {
	err := s.cacher.SetWithTTL(
		fmt.Sprintf("%s%d", baseKey, mnemonicId),
		LastSync{Timestamp: time.Now()},
		24*time.Hour,
	)
	if err != nil {
		return fmt.Errorf("failed to get last sync: %s", err)
	}

	return nil
}
