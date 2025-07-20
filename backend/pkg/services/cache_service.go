package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/doggyclub/backend/config"
	"github.com/doggyclub/backend/pkg/utils"
	"github.com/redis/go-redis/v9"
)

type CacheService struct {
	redis *redis.Client
	cfg   config.Config
	ctx   context.Context
}

func NewCacheService(redis *redis.Client, cfg config.Config) *CacheService {
	return &CacheService{
		redis: redis,
		cfg:   cfg,
		ctx:   context.Background(),
	}
}

// Cache keys constants
const (
	UserProfileKey        = "user:profile:%s"
	DogProfileKey         = "dog:profile:%s"
	PostKey               = "post:%s"
	TimelineKey           = "timeline:user:%s:page:%d"
	EncounterHistoryKey   = "encounters:user:%s:page:%d"
	GiftCatalogKey        = "gifts:catalog"
	GiftRankingsKey       = "gifts:rankings:%s"
	NotificationCountKey  = "notifications:unread:%s"
	LocationHashKey       = "location:hash:%s"
	ContentFilterKey      = "content:filters"
	SubscriptionPlansKey  = "subscription:plans"
	UserStatsKey          = "user:stats:%s"
	PopularPostsKey       = "posts:popular:%s" // period: daily, weekly, monthly
)

// Cache durations
const (
	ShortCacheDuration   = 5 * time.Minute
	MediumCacheDuration  = 30 * time.Minute
	LongCacheDuration    = 2 * time.Hour
	VeryLongCacheDuration = 24 * time.Hour
)

// Generic cache operations
func (s *CacheService) Set(key string, value interface{}, duration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return utils.WrapError(err, "failed to marshal cache data")
	}

	if err := s.redis.Set(s.ctx, key, data, duration).Err(); err != nil {
		return utils.WrapError(err, "failed to set cache")
	}

	return nil
}

func (s *CacheService) Get(key string, dest interface{}) error {
	data, err := s.redis.Get(s.ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return utils.NewAPIError("CACHE_MISS", "Cache key not found", nil)
		}
		return utils.WrapError(err, "failed to get cache")
	}

	if err := json.Unmarshal([]byte(data), dest); err != nil {
		return utils.WrapError(err, "failed to unmarshal cache data")
	}

	return nil
}

func (s *CacheService) Delete(key string) error {
	if err := s.redis.Del(s.ctx, key).Err(); err != nil {
		return utils.WrapError(err, "failed to delete cache key")
	}
	return nil
}

func (s *CacheService) DeletePattern(pattern string) error {
	keys, err := s.redis.Keys(s.ctx, pattern).Result()
	if err != nil {
		return utils.WrapError(err, "failed to get keys for pattern")
	}

	if len(keys) > 0 {
		if err := s.redis.Del(s.ctx, keys...).Err(); err != nil {
			return utils.WrapError(err, "failed to delete cache keys")
		}
	}

	return nil
}

func (s *CacheService) Exists(key string) (bool, error) {
	count, err := s.redis.Exists(s.ctx, key).Result()
	if err != nil {
		return false, utils.WrapError(err, "failed to check cache key existence")
	}
	return count > 0, nil
}

// Specific cache operations

// User profile caching
func (s *CacheService) CacheUserProfile(userID string, profile interface{}) error {
	key := fmt.Sprintf(UserProfileKey, userID)
	return s.Set(key, profile, MediumCacheDuration)
}

func (s *CacheService) GetUserProfile(userID string, dest interface{}) error {
	key := fmt.Sprintf(UserProfileKey, userID)
	return s.Get(key, dest)
}

func (s *CacheService) InvalidateUserProfile(userID string) error {
	key := fmt.Sprintf(UserProfileKey, userID)
	return s.Delete(key)
}

// Dog profile caching
func (s *CacheService) CacheDogProfile(dogID string, profile interface{}) error {
	key := fmt.Sprintf(DogProfileKey, dogID)
	return s.Set(key, profile, MediumCacheDuration)
}

func (s *CacheService) GetDogProfile(dogID string, dest interface{}) error {
	key := fmt.Sprintf(DogProfileKey, dogID)
	return s.Get(key, dest)
}

func (s *CacheService) InvalidateDogProfile(dogID string) error {
	key := fmt.Sprintf(DogProfileKey, dogID)
	return s.Delete(key)
}

// Post caching
func (s *CacheService) CachePost(postID string, post interface{}) error {
	key := fmt.Sprintf(PostKey, postID)
	return s.Set(key, post, ShortCacheDuration)
}

func (s *CacheService) GetPost(postID string, dest interface{}) error {
	key := fmt.Sprintf(PostKey, postID)
	return s.Get(key, dest)
}

func (s *CacheService) InvalidatePost(postID string) error {
	key := fmt.Sprintf(PostKey, postID)
	return s.Delete(key)
}

// Timeline caching
func (s *CacheService) CacheTimeline(userID string, page int, timeline interface{}) error {
	key := fmt.Sprintf(TimelineKey, userID, page)
	return s.Set(key, timeline, ShortCacheDuration)
}

func (s *CacheService) GetTimeline(userID string, page int, dest interface{}) error {
	key := fmt.Sprintf(TimelineKey, userID, page)
	return s.Get(key, dest)
}

func (s *CacheService) InvalidateUserTimeline(userID string) error {
	pattern := fmt.Sprintf("timeline:user:%s:page:*", userID)
	return s.DeletePattern(pattern)
}

// Gift catalog caching
func (s *CacheService) CacheGiftCatalog(catalog interface{}) error {
	return s.Set(GiftCatalogKey, catalog, LongCacheDuration)
}

func (s *CacheService) GetGiftCatalog(dest interface{}) error {
	return s.Get(GiftCatalogKey, dest)
}

func (s *CacheService) InvalidateGiftCatalog() error {
	return s.Delete(GiftCatalogKey)
}

// Gift rankings caching
func (s *CacheService) CacheGiftRankings(period string, rankings interface{}) error {
	key := fmt.Sprintf(GiftRankingsKey, period)
	return s.Set(key, rankings, MediumCacheDuration)
}

func (s *CacheService) GetGiftRankings(period string, dest interface{}) error {
	key := fmt.Sprintf(GiftRankingsKey, period)
	return s.Get(key, dest)
}

func (s *CacheService) InvalidateGiftRankings(period string) error {
	key := fmt.Sprintf(GiftRankingsKey, period)
	return s.Delete(key)
}

// Notification count caching
func (s *CacheService) CacheNotificationCount(userID string, count int64) error {
	key := fmt.Sprintf(NotificationCountKey, userID)
	return s.Set(key, count, ShortCacheDuration)
}

func (s *CacheService) GetNotificationCount(userID string) (int64, error) {
	key := fmt.Sprintf(NotificationCountKey, userID)
	var count int64
	err := s.Get(key, &count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *CacheService) InvalidateNotificationCount(userID string) error {
	key := fmt.Sprintf(NotificationCountKey, userID)
	return s.Delete(key)
}

// Location hash caching (for encounter detection)
func (s *CacheService) CacheLocationHash(hash string, userIDs []string) error {
	key := fmt.Sprintf(LocationHashKey, hash)
	return s.Set(key, userIDs, ShortCacheDuration)
}

func (s *CacheService) GetLocationHash(hash string) ([]string, error) {
	key := fmt.Sprintf(LocationHashKey, hash)
	var userIDs []string
	err := s.Get(key, &userIDs)
	if err != nil {
		return nil, err
	}
	return userIDs, nil
}

func (s *CacheService) InvalidateLocationHash(hash string) error {
	key := fmt.Sprintf(LocationHashKey, hash)
	return s.Delete(key)
}

// Content filters caching
func (s *CacheService) CacheContentFilters(filters interface{}) error {
	return s.Set(ContentFilterKey, filters, VeryLongCacheDuration)
}

func (s *CacheService) GetContentFilters(dest interface{}) error {
	return s.Get(ContentFilterKey, dest)
}

func (s *CacheService) InvalidateContentFilters() error {
	return s.Delete(ContentFilterKey)
}

// Subscription plans caching
func (s *CacheService) CacheSubscriptionPlans(plans interface{}) error {
	return s.Set(SubscriptionPlansKey, plans, VeryLongCacheDuration)
}

func (s *CacheService) GetSubscriptionPlans(dest interface{}) error {
	return s.Get(SubscriptionPlansKey, dest)
}

func (s *CacheService) InvalidateSubscriptionPlans() error {
	return s.Delete(SubscriptionPlansKey)
}

// User statistics caching
func (s *CacheService) CacheUserStats(userID string, stats interface{}) error {
	key := fmt.Sprintf(UserStatsKey, userID)
	return s.Set(key, stats, MediumCacheDuration)
}

func (s *CacheService) GetUserStats(userID string, dest interface{}) error {
	key := fmt.Sprintf(UserStatsKey, userID)
	return s.Get(key, dest)
}

func (s *CacheService) InvalidateUserStats(userID string) error {
	key := fmt.Sprintf(UserStatsKey, userID)
	return s.Delete(key)
}

// Popular posts caching
func (s *CacheService) CachePopularPosts(period string, posts interface{}) error {
	key := fmt.Sprintf(PopularPostsKey, period)
	return s.Set(key, posts, LongCacheDuration)
}

func (s *CacheService) GetPopularPosts(period string, dest interface{}) error {
	key := fmt.Sprintf(PopularPostsKey, period)
	return s.Get(key, dest)
}

func (s *CacheService) InvalidatePopularPosts(period string) error {
	key := fmt.Sprintf(PopularPostsKey, period)
	return s.Delete(key)
}

// Rate limiting using Redis
func (s *CacheService) CheckRateLimit(key string, limit int, window time.Duration) (bool, error) {
	pipe := s.redis.Pipeline()
	
	// Increment counter
	incr := pipe.Incr(s.ctx, key)
	
	// Set expiry if it's the first request
	pipe.Expire(s.ctx, key, window)
	
	_, err := pipe.Exec(s.ctx)
	if err != nil {
		return false, utils.WrapError(err, "failed to execute rate limit pipeline")
	}

	count, err := incr.Result()
	if err != nil {
		return false, utils.WrapError(err, "failed to get rate limit count")
	}

	return count <= int64(limit), nil
}

// Session management
func (s *CacheService) SetSession(sessionID string, data interface{}, duration time.Duration) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return s.Set(key, data, duration)
}

func (s *CacheService) GetSession(sessionID string, dest interface{}) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return s.Get(key, dest)
}

func (s *CacheService) DeleteSession(sessionID string) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return s.Delete(key)
}

// Batch operations
func (s *CacheService) SetMultiple(items map[string]interface{}, duration time.Duration) error {
	pipe := s.redis.Pipeline()
	
	for key, value := range items {
		data, err := json.Marshal(value)
		if err != nil {
			return utils.WrapError(err, "failed to marshal cache data for key: "+key)
		}
		pipe.Set(s.ctx, key, data, duration)
	}
	
	_, err := pipe.Exec(s.ctx)
	if err != nil {
		return utils.WrapError(err, "failed to execute batch set pipeline")
	}
	
	return nil
}

func (s *CacheService) GetMultiple(keys []string) (map[string]string, error) {
	if len(keys) == 0 {
		return make(map[string]string), nil
	}
	
	pipe := s.redis.Pipeline()
	cmds := make([]*redis.StringCmd, len(keys))
	
	for i, key := range keys {
		cmds[i] = pipe.Get(s.ctx, key)
	}
	
	_, err := pipe.Exec(s.ctx)
	if err != nil && err != redis.Nil {
		return nil, utils.WrapError(err, "failed to execute batch get pipeline")
	}
	
	result := make(map[string]string)
	for i, cmd := range cmds {
		val, err := cmd.Result()
		if err == redis.Nil {
			continue // Skip missing keys
		}
		if err != nil {
			return nil, utils.WrapError(err, "failed to get value for key: "+keys[i])
		}
		result[keys[i]] = val
	}
	
	return result, nil
}

// Cache warming (preload frequently accessed data)
func (s *CacheService) WarmCache() error {
	// This would be called on application startup to preload
	// frequently accessed data like gift catalog, subscription plans, etc.
	
	// Example implementation would fetch and cache:
	// - Gift catalog
	// - Subscription plans  
	// - Content filters
	// - Popular posts
	
	return nil
}

// Cache statistics
func (s *CacheService) GetCacheStats() (map[string]interface{}, error) {
	info, err := s.redis.Info(s.ctx, "memory").Result()
	if err != nil {
		return nil, utils.WrapError(err, "failed to get Redis memory info")
	}
	
	stats := map[string]interface{}{
		"redis_info": info,
	}
	
	return stats, nil
}