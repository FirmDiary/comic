package common

import (
	"comic/config"
	"github.com/silenceper/wechat/v2/cache"
)

func NewRedis() *cache.Redis {
	globalCfg := config.GetConfig()

	return cache.NewRedis(&cache.RedisOpts{
		Host: globalCfg.Redis.Host,
		Password: globalCfg.Redis.Password,
	})
}
