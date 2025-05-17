package config

import (
	"github.com/spf13/viper"
	"huinongfinancial/models"
)

func LoadConfig() {
	// 设置配置文件名
	viper.SetConfigName("config")
	// 设置配置文件类型
	viper.SetConfigType("yaml")
	// 设置配置文件路径
	viper.AddConfigPath(".")

	// 读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	// 将配置文件转换为结构体
	err = viper.Unmarshal(&models.Config)
	if err != nil {
		panic(err)
	}

	// 初始化配置
	InitConfig()
}

// InitConfig 初始化配置
func InitConfig() {
	// 如果 Server.Port 为空，则设置为 8080
	if models.Config.Server.Port == "" {
		models.Config.Server.Port = "8080"
	}

	// 如果 Server.ReadTimeout 为空，则设置为 10 秒
	if models.Config.Server.ReadTimeout == "" {
		models.Config.Server.ReadTimeout = "10s"
	}

	// 如果 Server.WriteTimeout 为空，则设置为 10 秒
	if models.Config.Server.WriteTimeout == "" {
		models.Config.Server.WriteTimeout = "10s"
	}
	
	// 如果 Database.Host 为空，则设置为 localhost
	if models.Config.Database.Host == "" {
		models.Config.Database.Host = "localhost"
	}
	
	// 如果 Database.Port 为空，则设置为 3306
	if models.Config.Database.Port == "" {
		models.Config.Database.Port = "3306"
	}

	// 如果 Database.User 为空，则设置为 root
	if models.Config.Database.User == "" {
		models.Config.Database.User = "root"
	}

	// 如果 Database.Password 为空，则设置为空
	if models.Config.Database.Password == "" {
		models.Config.Database.Password = ""
	}
}
