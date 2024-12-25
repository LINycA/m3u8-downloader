package dot

import (
	"strings"

	"github.com/spf13/viper"
)

// InitViper 初始化配置文件
func InitViper() {
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		Logger().WithError(err).Panic("配置文件导入")
		return
	}
	viper.SetDefault("semaphore", 6)
	viper.SetDefault("proxy", "direct")
	viper.SetDefault("cookie", "")
	viper.SetDefault("ua", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36")
	err = viper.WriteConfig()
	if err != nil {
		Logger().WithError(err).Error("初始化配置文件")
		return
	}
}
