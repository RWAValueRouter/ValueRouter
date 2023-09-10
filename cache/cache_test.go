package cache

import (
	"context"
	"github.com/RWAValueRouter/ValueRouter/common"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestRCache_SetGetValue(t *testing.T) {
	key := "first"
	if Cache.Client.Ping(context.Background()).String() != "ping: PONG" {
		return
	}
	err := Cache.SetValue(key, 100, 0)
	assert.NoError(t, err)
	ret, err := Cache.GetValue(key)
	if assert.NoError(t, err) {
		assert.Equal(t, "100", ret)
	}
}

func TestRCache_DeleteValue(t *testing.T) {
	if Cache.Client.Ping(context.Background()).String() != "ping: PONG" {
		return
	}
	key := "first"
	err := Cache.SetValue(key, 100, 0)
	assert.NoError(t, err)
	err = Cache.DeleteValue(key)
	assert.NoError(t, err)
	_, err = Cache.GetValue(key)
	if assert.Error(t, err) {
		assert.Equal(t, "redis: nil", err.Error())
	}
	err = Cache.DeleteValue("sec")
	assert.NoError(t, err)
}

func TestRCache_DeleteValueByPrefix(t *testing.T) {
	if Cache.Client.Ping(context.Background()).String() != "ping: PONG" {
		return
	}
	key := "first"
	for i := 0; i < 10; i++ {
		err := Cache.SetValue(key+strconv.Itoa(i), 100, 0)
		assert.NoError(t, err)
	}
	err := Cache.DeleteValueByPrefix(key)
	assert.NoError(t, err)
}

func TestRCache_SetJsonGetJsonValue(t *testing.T) {
	if Cache.Client.Ping(context.Background()).String() != "ping: PONG" {
		return
	}
	type people struct {
		Name string
		Age  int
	}
	p := people{
		Name: "clark",
		Age:  10,
	}
	err := Cache.SetJsonValue("people", &p, 0)
	assert.NoError(t, err)
	rp := &people{}
	err = Cache.GetJsonValue("people", rp)
	if assert.NoError(t, err) {
		assert.Equal(t, 10, rp.Age)
		assert.Equal(t, "clark", rp.Name)
	}
}

func init() {
	common.Conf.RedisConfig.Addr = "localhost:6379"
	common.Conf.RedisConfig.Password = ""
	common.Conf.RedisConfig.DB = 0
	common.Conf.RedisConfig.PoolSize = 100
	Init()
}
