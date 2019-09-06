package cache

import (
	"encoding/json"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	"live-service/app/util/config"
	"log"
)

// 获取缓存实例
func GetCache() cache.Cache {

	conf, err := config.GetAppConfig()
	if err != nil {
		log.Fatalf("app config error: %s", err)
	}

	cacheConfig := map[string]string{}

	switch conf.App.Cache.Driver {
	case "redis":
		cacheConfig["key"] = "cache"
		cacheConfig["conn"] = conf.App.Redis.Host+":"+conf.App.Redis.Port
		if conf.App.Redis.Password != "" {
			cacheConfig["password"] = conf.App.Redis.Password
		}
		break
	case "memcache":
		cacheConfig["conn"] = conf.App.MemCache.Host+":"+conf.App.MemCache.Port
		break
	case "file":
		cacheConfig["CachePath"] = "./storage/cache"
		cacheConfig["FileSuffix"] = ".cache"
		cacheConfig["DirectoryLevel"] = "2"
		cacheConfig["EmbedExpiry"] = "120"
		break
	default:
		log.Fatalf("cache driver error")
	}

	cacheConf, _ := json.Marshal(cacheConfig)

	bm, err := cache.NewCache(conf.App.Cache.Driver, string(cacheConf))

	if err != nil {
		log.Fatalf("cache config error: %s", err)
	}

	return bm
}
