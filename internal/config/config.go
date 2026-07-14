package config

import (
	"fmt"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	App     AppConfig     `mapstructure:"app"`
	Server  ServerConfig  `mapstructure:"server"`
	SQLite  SQLiteConfig  `mapstructure:"sqlite"`
	JWT     JWTConfig     `mapstructure:"jwt"`
	Payment PaymentConfig `mapstructure:"payment"`
}

type AppConfig struct {
	Name string `mapstructure:"name"`
	Env  string `mapstructure:"env"`
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type SQLiteConfig struct {
	Path string `mapstructure:"path"`
}

type JWTConfig struct {
	BuyerSecret string `mapstructure:"buyer_secret"`
	AdminSecret string `mapstructure:"admin_secret"`
	ExpireHours int    `mapstructure:"expire_hours"`
}

type PaymentConfig struct {
	Channels []ChannelConfig `mapstructure:"channels"`
}

type ChannelConfig struct {
	Name      string `mapstructure:"name"`
	Enabled   bool   `mapstructure:"enabled"`
	PID       string `mapstructure:"pid"`
	Key       string `mapstructure:"key"`
	APIURL    string `mapstructure:"api_url"`
	NotifyURL string `mapstructure:"notify_url"`
	ReturnURL string `mapstructure:"return_url"`
}

func Load(configPath string) (*Config, error) {
	// 自动加载 .env，忽略文件不存在的错误
	_ = godotenv.Load(".env")

	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")
	v.SetEnvPrefix("KDS")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &cfg, nil
}
