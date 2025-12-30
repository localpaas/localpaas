package realtimelog

import (
	"sync"
	"time"

	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
)

const (
	defaultLastActivePeriod = 10 * time.Second
	defaultExpiration       = 5 * time.Minute
)

type localStore struct {
	Store *Store
	Ts    time.Time
}

var (
	localStoreMap = make(map[string]*localStore)
	mu            sync.Mutex
)

func addLocalStore(key string, store *Store) {
	removeOldStores()
	mu.Lock()
	defer mu.Unlock()
	localStoreMap[key] = &localStore{Store: store, Ts: timeutil.NowUTC()}
}

func removeLocalStore(key string) {
	mu.Lock()
	defer mu.Unlock()
	delete(localStoreMap, key)
}

func removeOldStores() {
	mu.Lock()
	defer mu.Unlock()
	timeNow := timeutil.NowUTC()
	for key, localStore := range localStoreMap {
		if localStore.Ts.Add(defaultExpiration).Before(timeNow) {
			delete(localStoreMap, key)
		}
	}
}

func getLocalStore(key string, requiredActivePeriod time.Duration) *Store {
	mu.Lock()
	defer mu.Unlock()
	localStore, ok := localStoreMap[key]
	if !ok {
		return nil
	}
	if requiredActivePeriod == 0 {
		requiredActivePeriod = defaultLastActivePeriod
	}
	timeNow := timeutil.NowUTC()
	if requiredActivePeriod > 0 && localStore.Ts.Add(requiredActivePeriod).Before(timeNow) {
		delete(localStoreMap, key)
		return nil
	}
	localStore.Ts = timeNow
	return localStore.Store
}
