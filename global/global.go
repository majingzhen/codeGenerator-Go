package global

import (
	"fmt"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/mysql"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
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
		panic("db no before")
	}
	return db
}

// initDataBase 初始化db
func initDataBase() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=True&loc=Local",
		GVA_VP.GetString("database.username"),
		GVA_VP.GetString("database.password"),
		GVA_VP.GetString("database.host"),
		GVA_VP.GetString("database.port"),
		GVA_VP.GetString("database.db_name"))
	gcf := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   GVA_VP.GetString("database.table_prefix"), // 控制表前缀
			SingularTable: true,
		},
		Logger: logger.Default, // 控制是否sql输出，默认是不输出
	}
	if GVA_VP.GetBool("database.log_mode") {
		gcf.Logger = logger.Default.LogMode(logger.Info) // logger.Info 就会输出sql
	}
	GOrmDao, _ = gorm.Open(mysql.Open(dsn), gcf)
}

// initViper 初始化配置
func initViper() {
	GVA_VP = viper.New()
	GVA_VP.AddConfigPath(".")           // 添加配置文件搜索路径，点号为当前目录
	GVA_VP.AddConfigPath("./config")    // 添加多个搜索目录
	GVA_VP.SetConfigType("yml")         // 如果配置文件没有后缀，可以不用配置
	GVA_VP.SetConfigName("application") // 文件名，没有后缀
	// v.SetConfigFile("configs/app.yml")
	// 读取配置文件
	if err := GVA_VP.ReadInConfig(); err == nil {
		Logger.Info("use config file -> " + GVA_VP.ConfigFileUsed())
	}
}

// initLogger 初始化日志配置
func initLogger() {
	// 创建配置
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	// 初始化logger
	var err error
	Logger, err = config.Build()
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
	defer Logger.Sync() // 将缓冲区的日志刷新到输出
}

func InitConfig() {
	initLogger()
	initViper()
	if GVA_VP.GetString("gen_code.data_source") == "db" {
		initDataBase()
	}
}
