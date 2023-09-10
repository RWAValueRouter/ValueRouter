package cache

import (
	"context"
	"encoding/json"
	"github.com/RWAValueRouter/FastMulThreshold-DSA/log"
	"github.com/RWAValueRouter/ValueRouter/common"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	Cache = &RCache{}
)

type RCache struct {
	*redis.Client
	context.Context
}

func NewInstance() (err error) {
	Cache.Client = redis.NewClient(&redis.Options{
		Addr:     common.Conf.RedisConfig.Addr,
		Password: common.Conf.RedisConfig.Password,
		DB:       common.Conf.RedisConfig.DB,
		PoolSize: common.Conf.RedisConfig.PoolSize,
	})
	Cache.Context = context.Background()

	_, cancel := context.WithTimeout(Cache.Context, 5*time.Second)
	defer cancel()

	_, err = Cache.Ping(Cache.Context).Result()
	return err
}

func (r *RCache) SetValue(key string, value interface{}, expiration time.Duration) error {
	sc := r.Set(r.Context, key, value, expiration)
	if sc.Err() != nil {
		return sc.Err()
	}
	return nil
}

func (r *RCache) SetJsonValue(key string, value interface{}, expiration time.Duration) error {
	buf, err := json.Marshal(value)
	if err != nil {
		return err
	}
	sc := r.Set(r.Context, key, string(buf), expiration)
	if sc.Err() != nil {
		return sc.Err()
	}
	return nil
}

func (r *RCache) GetJsonValue(key string, ret interface{}) error {
	stringCmd := Cache.Get(Cache.Context, key)
	if stringCmd.Err() == nil {
		pstr, _ := stringCmd.Result()
		err := json.Unmarshal([]byte(pstr), ret)
		if err != nil {
			return err
		}
		return nil
	}
	return stringCmd.Err()
}

func (r *RCache) GetValue(key string) (string, error) {
	sc := r.Get(r.Context, key)
	if sc.Err() != nil {
		return "", sc.Err()
	}
	return sc.Result()
}

func (r *RCache) DeleteValue(key string) error {
	ic := r.Del(r.Context, key)
	if ic.Err() != nil {
		return ic.Err()
	}
	return nil
}

func (r *RCache) DeleteValueByPrefix(prefix string) error {
	iter := r.Scan(r.Context, 0, prefix+"*", 0).Iterator()
	for iter.Next(r.Context) {
		err := r.Del(r.Context, iter.Val()).Err()
		if err != nil {
			return err
		}
	}
	if err := iter.Err(); err != nil {
		return err

	}
	return nil
}

func Init() {
	err := NewInstance()
	if err != nil {
		log.Error("Connect Redis Error, Unable to Get redis instance")
	}
}
