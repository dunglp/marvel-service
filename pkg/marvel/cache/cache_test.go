package cache_test

import (
	"os"
	"testing"
	"time"
	mCache "xendit-technical-assessment/pkg/marvel/cache"

	"github.com/stretchr/testify/assert"

	"github.com/patrickmn/go-cache"
)

var cms mCache.CachedMarvelService

func TestMain(m *testing.M) {

	// prepare
	cms = mCache.NewCachedMarvelService(cache.New(5*time.Minute, 5*time.Minute))
	// run tests
	code := m.Run()

	// end
	os.Exit(code)
}

func TestCachedMarvelService_CacheCharacterIds(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		since := time.Now()
		cms.CacheCharacterIds([]int{1, 2, 3}, 3, since)
		found, cachedData := cms.CharacterIds()
		assert.Equal(t, true, found)
		assert.Equal(t, []int{1, 2, 3}, cachedData.IDs)
		assert.Equal(t, 3, cachedData.Total)
		assert.Equal(t, since, cachedData.Since)
	})
}

func TestCachedMarvelService_CharacterIds(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		since := time.Now()
		cms.CacheCharacterIds([]int{1, 2, 3}, 3, since)
		found, cachedData := cms.CharacterIds()
		assert.Equal(t, true, found)
		assert.Equal(t, []int{1, 2, 3}, cachedData.IDs)
		assert.Equal(t, since, cachedData.Since)
	})
}

func TestCachedMarvelService_UpdateCache(t *testing.T) {
	t.Run("success - modifiedIDs includes new ID only", func(t *testing.T) {
		since := time.Now()
		cms.CacheCharacterIds([]int{1, 2, 3}, 3, since)
		time.Sleep(5 * time.Second)

		since = time.Now()
		ids := cms.UpdateCache([]int{4}, since)

		found, cachedData := cms.CharacterIds()
		assert.Equal(t, true, found)
		assert.Equal(t, []int{1, 2, 3, 4}, cachedData.IDs)
		assert.Equal(t, []int{1, 2, 3, 4}, ids)
		assert.Equal(t, 4, cachedData.Total)
		assert.Equal(t, since, cachedData.Since)
	})
	t.Run("success - modifiedIDs includes existing IDs", func(t *testing.T) {
		since := time.Now()
		cms.CacheCharacterIds([]int{1, 2, 3}, 3, since)
		time.Sleep(5 * time.Second)

		since = time.Now()
		ids := cms.UpdateCache([]int{1, 2, 4}, since)

		found, cachedData := cms.CharacterIds()
		assert.Equal(t, true, found)
		assert.Equal(t, []int{1, 2, 3, 4}, cachedData.IDs)
		assert.Equal(t, []int{1, 2, 3, 4}, ids)
		assert.Equal(t, 4, cachedData.Total)
		assert.Equal(t, since, cachedData.Since)
	})
	t.Run("success - no new ID", func(t *testing.T) {
		since := time.Now()
		cms.CacheCharacterIds([]int{1, 2, 3}, 3, since)
		time.Sleep(5 * time.Second)

		since = time.Now()
		ids := cms.UpdateCache([]int{1, 2, 3}, since)

		found, cachedData := cms.CharacterIds()
		assert.Equal(t, true, found)
		assert.Equal(t, []int{1, 2, 3}, cachedData.IDs)
		assert.Equal(t, []int{1, 2, 3}, ids)
		assert.Equal(t, 3, cachedData.Total)
		assert.Equal(t, since, cachedData.Since)
	})
}
