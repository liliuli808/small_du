package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config 应用配置
type Config struct {
	App       AppConfig       `mapstructure:"app"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Redis     RedisConfig     `mapstructure:"redis"`
	AI        AIConfig        `mapstructure:"ai"`
	Bilibili  BilibiliConfig  `mapstructure:"bilibili"`
	Worker    WorkerConfig    `mapstructure:"worker"`
	RateLimit RateLimitConfig `mapstructure:"rate_limit"`
}

// AppConfig 应用配置
type AppConfig struct {
	Name string `mapstructure:"name"`
	Env  string `mapstructure:"env"`
	Port int    `mapstructure:"port"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
	MaxConns int    `mapstructure:"max_conns"`
}

// DSN 返回数据库连接字符串
func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.DBName, d.SSLMode)
}

// RedisConfig Redis配置
type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// AIConfig AI服务配置
type AIConfig struct {
	APIKey      string  `mapstructure:"api_key"`
	BaseURL     string  `mapstructure:"base_url"`
	Model       string  `mapstructure:"model"`
	MaxTokens   int     `mapstructure:"max_tokens"`
	Temperature float64 `mapstructure:"temperature"`
}

// BilibiliConfig B站配置
type BilibiliConfig struct {
	UserAgent string `mapstructure:"user_agent"`
	Cookie    string `mapstructure:"cookie"`
	Timeout   int    `mapstructure:"timeout"`
}

// WorkerConfig Worker配置
type WorkerConfig struct {
	Concurrency int `mapstructure:"concurrency"`
	Retry       int `mapstructure:"retry"`
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	IPLimitPerMinute    int `mapstructure:"ip_limit_per_minute"`
	VideoLimitPerMinute int `mapstructure:"video_limit_per_minute"`
}

// Load 加载配置
func Load(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	return &cfg, nil
}
