package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	Storage  StorageConfig
	JWT      JWTConfig
	APIKey   APIKeyConfig
	Log      LogConfig
	CORS     CORSConfig
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	Charset  string
}

type ServerConfig struct {
	Host      string
	Port      int
	APIPrefix string
}

type StorageConfig struct {
	Path        string
	MaxFileSize int64
}

type JWTConfig struct {
	Secret      string
	ExpireHours int
}

type APIKeyConfig struct {
	HashSalt string
}

type LogConfig struct {
	Level    string
	Format   string
	FilePath string // 日志文件路径，如果为空则只输出到 stdout
}

type CORSConfig struct {
	AllowOrigins []string
}

var AppConfig *Config

func Load() error {
	// 尝试加载 .env 文件（如果存在）
	_ = godotenv.Load()

	// 设置配置文件名称和路径
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./backend")

	// 尝试读取配置文件（如果存在）
	if err := viper.ReadInConfig(); err != nil {
		// 配置文件不存在不是错误，继续使用默认值
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// 其他错误才返回
			return fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// 设置默认值
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 3306)
	viper.SetDefault("database.user", "root")
	viper.SetDefault("database.password", "")
	viper.SetDefault("database.name", "dragonos_ci")
	viper.SetDefault("database.charset", "utf8mb4")

	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.api_prefix", "/api/v1")

	viper.SetDefault("storage.path", "./data/uploads")
	viper.SetDefault("storage.max_file_size", 104857600) // 100MB

	viper.SetDefault("jwt.secret", "change-me-in-production")
	viper.SetDefault("jwt.expire_hours", 24)

	viper.SetDefault("api_key.hash_salt", "change-me-in-production")

	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "json")
	viper.SetDefault("log.file_path", "") // 默认为空，只输出到 stdout

	viper.SetDefault("cors.allow_origins", []string{"http://localhost:3000", "http://localhost:5173"})

	// 从环境变量读取配置（环境变量优先级最高）
	viper.AutomaticEnv()

	// 绑定环境变量到配置键（支持两种格式）
	bindEnvVars()

	// Viper 会自动处理优先级：环境变量 > 配置文件 > 默认值
	// 同时支持环境变量格式（DB_HOST）和配置文件格式（database.host）
	AppConfig = &Config{
		Database: DatabaseConfig{
			Host:     getConfigValue("DB_HOST", "database.host", "localhost"),
			Port:     getConfigInt("DB_PORT", "database.port", 3306),
			User:     getConfigValue("DB_USER", "database.user", "root"),
			Password: getConfigValue("DB_PASSWORD", "database.password", ""),
			Name:     getConfigValue("DB_NAME", "database.name", "dragonos_ci"),
			Charset:  getConfigValue("DB_CHARSET", "database.charset", "utf8mb4"),
		},
		Server: ServerConfig{
			Host:      getConfigValue("SERVER_HOST", "server.host", "0.0.0.0"),
			Port:      getConfigInt("SERVER_PORT", "server.port", 8080),
			APIPrefix: getConfigValue("API_PREFIX", "server.api_prefix", "/api/v1"),
		},
		Storage: StorageConfig{
			Path:        getConfigValue("STORAGE_PATH", "storage.path", "./data/uploads"),
			MaxFileSize: getConfigInt64("MAX_FILE_SIZE", "storage.max_file_size", 104857600),
		},
		JWT: JWTConfig{
			Secret:      getConfigValue("JWT_SECRET", "jwt.secret", "change-me-in-production"),
			ExpireHours: getConfigInt("JWT_EXPIRE_HOURS", "jwt.expire_hours", 24),
		},
		APIKey: APIKeyConfig{
			HashSalt: getConfigValue("API_KEY_HASH_SALT", "api_key.hash_salt", "change-me-in-production"),
		},
		Log: LogConfig{
			Level:    getConfigValue("LOG_LEVEL", "log.level", "info"),
			Format:   getConfigValue("LOG_FORMAT", "log.format", "json"),
			FilePath: getConfigValue("LOG_FILE_PATH", "log.file_path", ""),
		},
		CORS: CORSConfig{
			AllowOrigins: getConfigStringSlice("CORS_ALLOW_ORIGINS", "cors.allow_origins", []string{"http://localhost:3000", "http://localhost:5173"}),
		},
	}

	// 确保存储目录存在
	if err := os.MkdirAll(AppConfig.Storage.Path, 0755); err != nil {
		return fmt.Errorf("failed to create storage directory: %w", err)
	}

	return nil
}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.Name, c.Charset)
}

func (c *JWTConfig) ExpireDuration() time.Duration {
	return time.Duration(c.ExpireHours) * time.Hour
}

// bindEnvVars 绑定环境变量到配置键
func bindEnvVars() {
	// 数据库配置
	viper.BindEnv("DB_HOST", "DB_HOST")
	viper.BindEnv("DB_PORT", "DB_PORT")
	viper.BindEnv("DB_USER", "DB_USER")
	viper.BindEnv("DB_PASSWORD", "DB_PASSWORD")
	viper.BindEnv("DB_NAME", "DB_NAME")
	viper.BindEnv("DB_CHARSET", "DB_CHARSET")

	// 服务器配置
	viper.BindEnv("SERVER_HOST", "SERVER_HOST")
	viper.BindEnv("SERVER_PORT", "SERVER_PORT")
	viper.BindEnv("API_PREFIX", "API_PREFIX")

	// 存储配置
	viper.BindEnv("STORAGE_PATH", "STORAGE_PATH")
	viper.BindEnv("MAX_FILE_SIZE", "MAX_FILE_SIZE")

	// JWT配置
	viper.BindEnv("JWT_SECRET", "JWT_SECRET")
	viper.BindEnv("JWT_EXPIRE_HOURS", "JWT_EXPIRE_HOURS")

	// API Key配置
	viper.BindEnv("API_KEY_HASH_SALT", "API_KEY_HASH_SALT")

	// 日志配置
	viper.BindEnv("LOG_LEVEL", "LOG_LEVEL")
	viper.BindEnv("LOG_FORMAT", "LOG_FORMAT")
	viper.BindEnv("LOG_FILE_PATH", "LOG_FILE_PATH")

	// CORS配置
	viper.BindEnv("CORS_ALLOW_ORIGINS", "CORS_ALLOW_ORIGINS")
}

// getConfigValue 获取配置值，优先级：环境变量 > 配置文件 > 默认值
func getConfigValue(envKey, configKey, defaultValue string) string {
	if viper.IsSet(envKey) {
		return viper.GetString(envKey)
	}
	if viper.IsSet(configKey) {
		return viper.GetString(configKey)
	}
	return defaultValue
}

// getConfigInt 获取整数配置值
func getConfigInt(envKey, configKey string, defaultValue int) int {
	if viper.IsSet(envKey) {
		return viper.GetInt(envKey)
	}
	if viper.IsSet(configKey) {
		return viper.GetInt(configKey)
	}
	return defaultValue
}

// getConfigInt64 获取int64配置值
func getConfigInt64(envKey, configKey string, defaultValue int64) int64 {
	if viper.IsSet(envKey) {
		return viper.GetInt64(envKey)
	}
	if viper.IsSet(configKey) {
		return viper.GetInt64(configKey)
	}
	return defaultValue
}

// getConfigStringSlice 获取字符串切片配置值
func getConfigStringSlice(envKey, configKey string, defaultValue []string) []string {
	if viper.IsSet(envKey) {
		return viper.GetStringSlice(envKey)
	}
	if viper.IsSet(configKey) {
		return viper.GetStringSlice(configKey)
	}
	return defaultValue
}
