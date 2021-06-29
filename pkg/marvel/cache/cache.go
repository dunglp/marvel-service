package cache

import (
	"time"

	"github.com/rs/zerolog/log"

	"github.com/patrickmn/go-cache"
)

const (
	characterIds = "character-ids"
)

type CharacterIDs struct {
	IDs   []int
	Total int
	Since time.Time
}

type CachedMarvelService struct {
	cache *cache.Cache
}

func NewCachedMarvelService(cache *cache.Cache) CachedMarvelService {
	return CachedMarvelService{
		cache,
	}
}

func (cms CachedMarvelService) CharacterIds() (found bool, data CharacterIDs) {
	cached, found := cms.cache.Get(characterIds)
	if found {
		return found, cached.(CharacterIDs)
	}

	return
}

func (cms CachedMarvelService) CacheCharacterIds(ids []int, total int, since time.Time) {
	log.Info().Msgf("CacheCharacterIds. Total %d at : %s", total, since.Format(time.RFC3339))
	cms.cache.Set(characterIds,
		CharacterIDs{
			ids,
			total,
			since,
		},
		cache.NoExpiration)
}

func (cms CachedMarvelService) UpdateCache(modifiedIds []int, since time.Time) (ids []int) {
	log.Info().Msgf("UpdateCache with %d characters at : %s", len(modifiedIds), since.Format(time.RFC3339))

	cached, _ := cms.cache.Get(characterIds)
	cachedData := cached.(CharacterIDs)

	ids = cachedData.IDs

	if len(modifiedIds) > 0 {
		ids = AddNewIds(ids, modifiedIds)
	}
	cms.CacheCharacterIds(ids, len(ids), since)
	return
}

func AddNewIds(ids []int, mIds []int) []int {
	m := make(map[int]int)
	for _, id := range ids {
		m[id] = id
	}

	for _, id := range mIds {
		if _, found := m[id]; !found {
			ids = append(ids, id)
		}
	}

	return ids
}