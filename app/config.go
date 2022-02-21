package app

import "github.com/spf13/viper"

var conf *viper.Viper

// LoadConfig 設定ファイルを読み込む
func LoadConfig(env string) error {
	conf = viper.New()
	conf.SetConfigName(env)
	conf.SetConfigType("toml")
	conf.AddConfigPath("./configs/")
	conf.AddConfigPath("/usr/local/bin/configs/")
	//conf.SetConfigFile(filepath)
	return conf.ReadInConfig()
}

// Config 設定を取得
func Config() *viper.Viper {
	return conf
}
