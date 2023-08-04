package global

import (
	"sync"

	//"zeng-frame-v1/frame/utils/timer"
	"github.com/songzhibin97/gkit/cache/local_cache"

	"golang.org/x/sync/singleflight"

	"go.uber.org/zap"

	//"zeng-frame-v1/frame/config"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var (
	GOrmDao    *gorm.DB
	GVA_DBList map[string]*gorm.DB
	GVA_REDIS  *redis.Client
	//GVA_CONFIG config.Server
	GVA_VP *viper.Viper
	// Logger    *oplogging.Logger
	Logger *zap.Logger
	//GVA_Timer               timer.Timer = timer.NewTimerTask()
	GVA_Concurrency_Control = &singleflight.Group{}

	BlackCache local_cache.Cache
	lock       sync.RWMutex
)

// GetGlobalDBByDBName 通过名称获取db list中的db
func GetGlobalDBByDBName(dbname string) *gorm.DB {
	lock.RLock()
	defer lock.RUnlock()
	return GVA_DBList[dbname]
}

// MustGetGlobalDBByDBName 通过名称获取db 如果不存在则panic
func MustGetGlobalDBByDBName(dbname string) *gorm.DB {
	lock.RLock()
	defer lock.RUnlock()
	db, ok := GVA_DBList[dbname]
	if !ok || db == nil {
		panic("db no init")
	}
	return db
}
