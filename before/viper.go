package before

import (
	"codeGenerator-Go/global"
	"github.com/spf13/viper"
	"log"
)

// Viper 初始化配置
func Viper() {
	global.GVA_VP = viper.New()
	global.GVA_VP.AddConfigPath(".")        // 添加配置文件搜索路径，点号为当前目录
	global.GVA_VP.AddConfigPath("./config") // 添加多个搜索目录
	global.GVA_VP.SetConfigType("yaml")     // 如果配置文件没有后缀，可以不用配置
	global.GVA_VP.SetConfigName("config")   // 文件名，没有后缀
	// v.SetConfigFile("configs/app.yml")
	// 读取配置文件
	if err := global.GVA_VP.ReadInConfig(); err == nil {
		log.Printf("use config file -> %s\n", global.GVA_VP.ConfigFileUsed())
	}
}
