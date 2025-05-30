package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisClient Redis客户端封装
type RedisClient struct {
	client *redis.Client
}

// CacheInterface 缓存接口
type CacheInterface interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	GetJSON(ctx context.Context, key string, dest interface{}) error
	SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	Expire(ctx context.Context, key string, expiration time.Duration) error
	Incr(ctx context.Context, key string) (int64, error)
	Decr(ctx context.Context, key string) (int64, error)
	IncrBy(ctx context.Context, key string, value int64) (int64, error)
	DecrBy(ctx context.Context, key string, value int64) (int64, error)
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error)
	GetPattern(ctx context.Context, pattern string) ([]string, error)
	DeletePattern(ctx context.Context, pattern string) error
	HSet(ctx context.Context, key string, field string, value interface{}) error
	HGet(ctx context.Context, key string, field string) (string, error)
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	HDel(ctx context.Context, key string, fields ...string) error
	SAdd(ctx context.Context, key string, members ...interface{}) error
	SMembers(ctx context.Context, key string) ([]string, error)
	SRem(ctx context.Context, key string, members ...interface{}) error
	ZAdd(ctx context.Context, key string, score float64, member interface{}) error
	ZRange(ctx context.Context, key string, start, stop int64) ([]string, error)
	ZRangeWithScores(ctx context.Context, key string, start, stop int64) ([]redis.Z, error)
	ZRem(ctx context.Context, key string, members ...interface{}) error
	Publish(ctx context.Context, channel string, message interface{}) error
	Subscribe(ctx context.Context, channel string) (<-chan string, error)
}

// NewRedisClient 创建Redis客户端
func NewRedisClient(addr, password string, db int) CacheInterface {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
		PoolSize: 10,
	})

	return &RedisClient{
		client: rdb,
	}
}

// Set 设置键值对
func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

// Get 获取值
func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	result, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key %s not found", key)
	}
	return result, err
}

// GetJSON 获取JSON值并反序列化
func (r *RedisClient) GetJSON(ctx context.Context, key string, dest interface{}) error {
	result, err := r.Get(ctx, key)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(result), dest)
}

// SetJSON 序列化JSON并设置
func (r *RedisClient) SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return r.Set(ctx, key, string(jsonData), expiration)
}

// Delete 删除键
func (r *RedisClient) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// Exists 检查键是否存在
func (r *RedisClient) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.client.Exists(ctx, key).Result()
	return result > 0, err
}

// Expire 设置过期时间
func (r *RedisClient) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return r.client.Expire(ctx, key, expiration).Err()
}

// Incr 自增
func (r *RedisClient) Incr(ctx context.Context, key string) (int64, error) {
	return r.client.Incr(ctx, key).Result()
}

// Decr 自减
func (r *RedisClient) Decr(ctx context.Context, key string) (int64, error) {
	return r.client.Decr(ctx, key).Result()
}

// IncrBy 按指定值自增
func (r *RedisClient) IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	return r.client.IncrBy(ctx, key, value).Result()
}

// DecrBy 按指定值自减
func (r *RedisClient) DecrBy(ctx context.Context, key string, value int64) (int64, error) {
	return r.client.DecrBy(ctx, key, value).Result()
}

// SetNX 设置键值对（仅当键不存在时）
func (r *RedisClient) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	return r.client.SetNX(ctx, key, value, expiration).Result()
}

// GetPattern 根据模式获取键列表
func (r *RedisClient) GetPattern(ctx context.Context, pattern string) ([]string, error) {
	return r.client.Keys(ctx, pattern).Result()
}

// DeletePattern 根据模式删除键
func (r *RedisClient) DeletePattern(ctx context.Context, pattern string) error {
	keys, err := r.GetPattern(ctx, pattern)
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		return nil
	}
	return r.client.Del(ctx, keys...).Err()
}

// HSet 设置哈希字段
func (r *RedisClient) HSet(ctx context.Context, key string, field string, value interface{}) error {
	return r.client.HSet(ctx, key, field, value).Err()
}

// HGet 获取哈希字段值
func (r *RedisClient) HGet(ctx context.Context, key string, field string) (string, error) {
	result, err := r.client.HGet(ctx, key, field).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("field %s in key %s not found", field, key)
	}
	return result, err
}

// HGetAll 获取哈希所有字段
func (r *RedisClient) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return r.client.HGetAll(ctx, key).Result()
}

// HDel 删除哈希字段
func (r *RedisClient) HDel(ctx context.Context, key string, fields ...string) error {
	return r.client.HDel(ctx, key, fields...).Err()
}

// SAdd 向集合添加成员
func (r *RedisClient) SAdd(ctx context.Context, key string, members ...interface{}) error {
	return r.client.SAdd(ctx, key, members...).Err()
}

// SMembers 获取集合所有成员
func (r *RedisClient) SMembers(ctx context.Context, key string) ([]string, error) {
	return r.client.SMembers(ctx, key).Result()
}

// SRem 从集合移除成员
func (r *RedisClient) SRem(ctx context.Context, key string, members ...interface{}) error {
	return r.client.SRem(ctx, key, members...).Err()
}

// ZAdd 向有序集合添加成员
func (r *RedisClient) ZAdd(ctx context.Context, key string, score float64, member interface{}) error {
	return r.client.ZAdd(ctx, key, redis.Z{Score: score, Member: member}).Err()
}

// ZRange 获取有序集合范围内的成员
func (r *RedisClient) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return r.client.ZRange(ctx, key, start, stop).Result()
}

// ZRangeWithScores 获取有序集合范围内的成员及分数
func (r *RedisClient) ZRangeWithScores(ctx context.Context, key string, start, stop int64) ([]redis.Z, error) {
	return r.client.ZRangeWithScores(ctx, key, start, stop).Result()
}

// ZRem 从有序集合移除成员
func (r *RedisClient) ZRem(ctx context.Context, key string, members ...interface{}) error {
	return r.client.ZRem(ctx, key, members...).Err()
}

// Publish 发布消息到频道
func (r *RedisClient) Publish(ctx context.Context, channel string, message interface{}) error {
	return r.client.Publish(ctx, channel, message).Err()
}

// Subscribe 订阅频道
func (r *RedisClient) Subscribe(ctx context.Context, channel string) (<-chan string, error) {
	pubsub := r.client.Subscribe(ctx, channel)
	ch := make(chan string, 100)

	go func() {
		defer close(ch)
		defer pubsub.Close()

		for msg := range pubsub.Channel() {
			ch <- msg.Payload
		}
	}()

	return ch, nil
}

// CacheService 缓存服务封装
type CacheService struct {
	client CacheInterface
}

// NewCacheService 创建缓存服务
func NewCacheService(client CacheInterface) *CacheService {
	return &CacheService{
		client: client,
	}
}

// 用户相关缓存方法

// SetUserCache 设置用户缓存
func (c *CacheService) SetUserCache(ctx context.Context, userID uint64, userData interface{}) error {
	key := fmt.Sprintf("user:%d", userID)
	return c.client.SetJSON(ctx, key, userData, 30*time.Minute)
}

// GetUserCache 获取用户缓存
func (c *CacheService) GetUserCache(ctx context.Context, userID uint64, dest interface{}) error {
	key := fmt.Sprintf("user:%d", userID)
	return c.client.GetJSON(ctx, key, dest)
}

// DeleteUserCache 删除用户缓存
func (c *CacheService) DeleteUserCache(ctx context.Context, userID uint64) error {
	key := fmt.Sprintf("user:%d", userID)
	return c.client.Delete(ctx, key)
}

// SetSessionCache 设置会话缓存
func (c *CacheService) SetSessionCache(ctx context.Context, sessionID string, sessionData interface{}) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return c.client.SetJSON(ctx, key, sessionData, 24*time.Hour)
}

// GetSessionCache 获取会话缓存
func (c *CacheService) GetSessionCache(ctx context.Context, sessionID string, dest interface{}) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return c.client.GetJSON(ctx, key, dest)
}

// DeleteSessionCache 删除会话缓存
func (c *CacheService) DeleteSessionCache(ctx context.Context, sessionID string) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return c.client.Delete(ctx, key)
}

// SetSMSCode 设置短信验证码
func (c *CacheService) SetSMSCode(ctx context.Context, phone, code string) error {
	key := fmt.Sprintf("sms_code:%s", phone)
	return c.client.Set(ctx, key, code, 5*time.Minute)
}

// GetSMSCode 获取短信验证码
func (c *CacheService) GetSMSCode(ctx context.Context, phone string) (string, error) {
	key := fmt.Sprintf("sms_code:%s", phone)
	return c.client.Get(ctx, key)
}

// DeleteSMSCode 删除短信验证码
func (c *CacheService) DeleteSMSCode(ctx context.Context, phone string) error {
	key := fmt.Sprintf("sms_code:%s", phone)
	return c.client.Delete(ctx, key)
}

// SetLimitCounter 设置限流计数器
func (c *CacheService) SetLimitCounter(ctx context.Context, key string, count int64, expiration time.Duration) error {
	return c.client.Set(ctx, fmt.Sprintf("limit:%s", key), count, expiration)
}

// IncrLimitCounter 增加限流计数器
func (c *CacheService) IncrLimitCounter(ctx context.Context, key string) (int64, error) {
	return c.client.Incr(ctx, fmt.Sprintf("limit:%s", key))
}

// GetLimitCounter 获取限流计数器
func (c *CacheService) GetLimitCounter(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, fmt.Sprintf("limit:%s", key))
}

// 文章相关缓存

// SetArticleCache 设置文章缓存
func (c *CacheService) SetArticleCache(ctx context.Context, articleID uint64, articleData interface{}) error {
	key := fmt.Sprintf("article:%d", articleID)
	return c.client.SetJSON(ctx, key, articleData, 1*time.Hour)
}

// GetArticleCache 获取文章缓存
func (c *CacheService) GetArticleCache(ctx context.Context, articleID uint64, dest interface{}) error {
	key := fmt.Sprintf("article:%d", articleID)
	return c.client.GetJSON(ctx, key, dest)
}

// SetArticleListCache 设置文章列表缓存
func (c *CacheService) SetArticleListCache(ctx context.Context, cacheKey string, articles interface{}) error {
	key := fmt.Sprintf("article_list:%s", cacheKey)
	return c.client.SetJSON(ctx, key, articles, 15*time.Minute)
}

// GetArticleListCache 获取文章列表缓存
func (c *CacheService) GetArticleListCache(ctx context.Context, cacheKey string, dest interface{}) error {
	key := fmt.Sprintf("article_list:%s", cacheKey)
	return c.client.GetJSON(ctx, key, dest)
}

// 贷款相关缓存

// SetLoanProductCache 设置贷款产品缓存
func (c *CacheService) SetLoanProductCache(ctx context.Context, productID uint64, productData interface{}) error {
	key := fmt.Sprintf("loan_product:%d", productID)
	return c.client.SetJSON(ctx, key, productData, 2*time.Hour)
}

// GetLoanProductCache 获取贷款产品缓存
func (c *CacheService) GetLoanProductCache(ctx context.Context, productID uint64, dest interface{}) error {
	key := fmt.Sprintf("loan_product:%d", productID)
	return c.client.GetJSON(ctx, key, dest)
}

// InvalidateLoanProductCache 失效贷款产品缓存
func (c *CacheService) InvalidateLoanProductCache(ctx context.Context, productID uint64) error {
	key := fmt.Sprintf("loan_product:%d", productID)
	return c.client.Delete(ctx, key)
}

// 农机相关缓存

// SetMachineCache 设置农机缓存
func (c *CacheService) SetMachineCache(ctx context.Context, machineID uint64, machineData interface{}) error {
	key := fmt.Sprintf("machine:%d", machineID)
	return c.client.SetJSON(ctx, key, machineData, 1*time.Hour)
}

// GetMachineCache 获取农机缓存
func (c *CacheService) GetMachineCache(ctx context.Context, machineID uint64, dest interface{}) error {
	key := fmt.Sprintf("machine:%d", machineID)
	return c.client.GetJSON(ctx, key, dest)
}

// 系统配置缓存

// SetConfigCache 设置系统配置缓存
func (c *CacheService) SetConfigCache(ctx context.Context, configKey string, configValue interface{}) error {
	key := fmt.Sprintf("config:%s", configKey)
	return c.client.SetJSON(ctx, key, configValue, 24*time.Hour)
}

// GetConfigCache 获取系统配置缓存
func (c *CacheService) GetConfigCache(ctx context.Context, configKey string, dest interface{}) error {
	key := fmt.Sprintf("config:%s", configKey)
	return c.client.GetJSON(ctx, key, dest)
}

// InvalidateConfigCache 失效系统配置缓存
func (c *CacheService) InvalidateConfigCache(ctx context.Context, configKey string) error {
	key := fmt.Sprintf("config:%s", configKey)
	return c.client.Delete(ctx, key)
}

// 分布式锁

// Lock 获取分布式锁
func (c *CacheService) Lock(ctx context.Context, lockKey string, expiration time.Duration) (bool, error) {
	key := fmt.Sprintf("lock:%s", lockKey)
	return c.client.SetNX(ctx, key, "locked", expiration)
}

// Unlock 释放分布式锁
func (c *CacheService) Unlock(ctx context.Context, lockKey string) error {
	key := fmt.Sprintf("lock:%s", lockKey)
	return c.client.Delete(ctx, key)
}

// IsLocked 检查是否已锁定
func (c *CacheService) IsLocked(ctx context.Context, lockKey string) (bool, error) {
	key := fmt.Sprintf("lock:%s", lockKey)
	return c.client.Exists(ctx, key)
}

// Publish 发布消息到频道
func (c *CacheService) Publish(ctx context.Context, channel string, message interface{}) error {
	return c.client.Publish(ctx, channel, message)
}

// Subscribe 订阅频道
func (c *CacheService) Subscribe(ctx context.Context, channel string) (<-chan string, error) {
	return c.client.Subscribe(ctx, channel)
}
