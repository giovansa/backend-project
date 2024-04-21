package internal

import "github.com/spf13/viper"

type Config struct {
	App AppConfig `mapstructure:"app"`
	DB  Database  `mapstructure:"database"`
}
type AppConfig struct {
	Env string `mapstructure:"env"`
}

type Database struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DatabaseName string `mapstructure:"database"`
}

func LoadConfig(path string) (Config, error) {
	var config Config
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		return config, err
	}
	if err = viper.Unmarshal(&config); err != nil {
		return config, err
	}
	return config, nil
}
