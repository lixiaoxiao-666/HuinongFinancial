package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config 应用配置
type Config struct {
	App       AppConfig       `mapstructure:"app"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Redis     RedisConfig     `mapstructure:"redis"`
	JWT       JWTConfig       `mapstructure:"jwt"`
	Session   SessionConfig   `mapstructure:"session"`
	Dify      DifyConfig      `mapstructure:"dify"`
	SMS       SMSConfig       `mapstructure:"sms"`
	File      FileConfig      `mapstructure:"file"`
	Log       LogConfig       `mapstructure:"log"`
	CORS      CORSConfig      `mapstructure:"cors"`
	RateLimit RateLimitConfig `mapstructure:"rate_limit"`
	Offline   OfflineConfig   `mapstructure:"offline"`
}

// AppConfig 应用基础配置
type AppConfig struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
	Env     string `mapstructure:"env"`
	Host    string `mapstructure:"host"`
	Port    int    `mapstructure:"port"`
	Mode    string `mapstructure:"mode"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver          string `mapstructure:"driver"`
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	Database        string `mapstructure:"database"`
	Charset         string `mapstructure:"charset"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Password     string `mapstructure:"password"`
	Database     int    `mapstructure:"database"`
	PoolSize     int    `mapstructure:"pool_size"`
	MinIdleConns int    `mapstructure:"min_idle_conns"`
}

// JWTConfig JWT认证配置
type JWTConfig struct {
	SecretKey        string `mapstructure:"secret_key"`
	ExpiresIn        int64  `mapstructure:"expires_in"`
	RefreshExpiresIn int64  `mapstructure:"refresh_expires_in"`
}

// SessionConfig 会话管理配置
type SessionConfig struct {
	Redis    SessionRedisConfig    `mapstructure:"redis"`
	Settings SessionSettingsConfig `mapstructure:"settings"`
	Security SessionSecurityConfig `mapstructure:"security"`
}

// SessionRedisConfig 会话Redis配置
type SessionRedisConfig struct {
	DB           int           `mapstructure:"db"`
	PoolSize     int           `mapstructure:"pool_size"`
	MinIdleConns int           `mapstructure:"min_idle_conns"`
	MaxRetries   int           `mapstructure:"max_retries"`
	DialTimeout  time.Duration `mapstructure:"dial_timeout"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

// SessionSettingsConfig 会话设置配置
type SessionSettingsConfig struct {
	AccessTokenTTL     time.Duration `mapstructure:"access_token_ttl"`
	RefreshTokenTTL    time.Duration `mapstructure:"refresh_token_ttl"`
	MaxSessionsPerUser int           `mapstructure:"max_sessions_per_user"`
	SessionTimeout     time.Duration `mapstructure:"session_timeout"`
	SingleDeviceLogin  bool          `mapstructure:"single_device_login"`
	CleanupInterval    time.Duration `mapstructure:"cleanup_interval"`
	BatchCleanupSize   int           `mapstructure:"batch_cleanup_size"`
}

// SessionSecurityConfig 会话安全配置
type SessionSecurityConfig struct {
	ValidateIP            bool    `mapstructure:"validate_ip"`
	AllowIPChange         bool    `mapstructure:"allow_ip_change"`
	ValidateDeviceID      bool    `mapstructure:"validate_device_id"`
	AllowDeviceChange     bool    `mapstructure:"allow_device_change"`
	MaxConcurrentSessions int     `mapstructure:"max_concurrent_sessions"`
	KickOldestSession     bool    `mapstructure:"kick_oldest_session"`
	EnableRiskDetection   bool    `mapstructure:"enable_risk_detection"`
	RiskThreshold         int     `mapstructure:"risk_threshold"`
	AutoRefresh           bool    `mapstructure:"auto_refresh"`
	RefreshThreshold      float64 `mapstructure:"refresh_threshold"`
}

// DifyConfig Dify AI集成配置
type DifyConfig struct {
	APIURL     string            `mapstructure:"api_url"`
	APIKey     string            `mapstructure:"api_key"`
	APIToken   string            `mapstructure:"api_token"`
	Timeout    int               `mapstructure:"timeout"`
	RetryTimes int               `mapstructure:"retry_times"`
	Workflows  map[string]string `mapstructure:"workflows"`
}

// SMSConfig 短信服务配置
type SMSConfig struct {
	Provider      string            `mapstructure:"provider"`
	AccessKey     string            `mapstructure:"access_key"`
	AccessSecret  string            `mapstructure:"access_secret"`
	SignName      string            `mapstructure:"sign_name"`
	TemplateCodes map[string]string `mapstructure:"template_codes"`
}

// FileConfig 文件上传配置
type FileConfig struct {
	StorageType  string   `mapstructure:"storage_type"`
	UploadPath   string   `mapstructure:"upload_path"`
	MaxFileSize  int64    `mapstructure:"max_file_size"`
	AllowedTypes []string `mapstructure:"allowed_types"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	Output     string `mapstructure:"output"`
	FilePath   string `mapstructure:"file_path"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

// CORSConfig CORS跨域配置
type CORSConfig struct {
	AllowOrigins     []string `mapstructure:"allow_origins"`
	AllowMethods     []string `mapstructure:"allow_methods"`
	AllowHeaders     []string `mapstructure:"allow_headers"`
	ExposeHeaders    []string `mapstructure:"expose_headers"`
	AllowCredentials bool     `mapstructure:"allow_credentials"`
	MaxAge           int      `mapstructure:"max_age"`
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	Enabled           bool `mapstructure:"enabled"`
	RequestsPerMinute int  `mapstructure:"requests_per_minute"`
	Burst             int  `mapstructure:"burst"`
}

// OfflineConfig 离线功能配置
type OfflineConfig struct {
	Enabled      bool `mapstructure:"enabled"`
	SyncInterval int  `mapstructure:"sync_interval"`
	QueueSize    int  `mapstructure:"queue_size"`
	RetryTimes   int  `mapstructure:"retry_times"`
}

// LoadConfig 加载配置文件
func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)

	// 设置环境变量前缀
	viper.SetEnvPrefix("HUINONG")
	viper.AutomaticEnv()

	// 设置环境变量分隔符
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 验证配置
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("配置验证失败: %w", err)
	}

	// 确保必要的目录存在
	if err := cfg.ensureDirectories(); err != nil {
		return nil, fmt.Errorf("创建必要目录失败: %w", err)
	}

	return &cfg, nil
}

// validate 验证配置
func (c *Config) validate() error {
	// 验证应用配置
	if c.App.Name == "" {
		return fmt.Errorf("应用名称不能为空")
	}
	if c.App.Port <= 0 || c.App.Port > 65535 {
		return fmt.Errorf("端口号必须在1-65535之间")
	}

	// 验证数据库配置
	if c.Database.Host == "" {
		return fmt.Errorf("数据库主机不能为空")
	}
	if c.Database.Username == "" {
		return fmt.Errorf("数据库用户名不能为空")
	}
	if c.Database.Database == "" {
		return fmt.Errorf("数据库名不能为空")
	}

	// 验证Redis配置
	if c.Redis.Host == "" {
		return fmt.Errorf("Redis主机不能为空")
	}

	// 验证JWT配置
	if c.JWT.SecretKey == "" {
		return fmt.Errorf("JWT密钥不能为空")
	}
	if len(c.JWT.SecretKey) < 32 {
		return fmt.Errorf("JWT密钥长度不能少于32位")
	}

	return nil
}

// ensureDirectories 确保必要的目录存在
func (c *Config) ensureDirectories() error {
	dirs := []string{
		c.File.UploadPath,
	}

	// 如果日志输出到文件，创建日志目录
	if c.Log.Output == "file" || c.Log.Output == "both" {
		dirs = append(dirs, getLogDir(c.Log.FilePath))
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("创建目录 %s 失败: %w", dir, err)
		}
	}

	return nil
}

// getLogDir 从日志文件路径获取目录
func getLogDir(filePath string) string {
	if filePath == "" {
		return "./logs"
	}

	// 获取目录部分
	lastSlash := strings.LastIndex(filePath, "/")
	if lastSlash == -1 {
		return "."
	}

	return filePath[:lastSlash]
}

// GetDSN 获取数据库连接字符串
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		c.Username, c.Password, c.Host, c.Port, c.Database, c.Charset)
}

// GetRedisAddr 获取Redis连接地址
func (c *RedisConfig) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// GetJWTExpirationDuration 获取JWT过期时间
func (c *JWTConfig) GetJWTExpirationDuration() time.Duration {
	return time.Duration(c.ExpiresIn) * time.Second
}

// GetRefreshExpirationDuration 获取刷新令牌过期时间
func (c *JWTConfig) GetRefreshExpirationDuration() time.Duration {
	return time.Duration(c.RefreshExpiresIn) * time.Second
}

// IsDevelopment 判断是否为开发环境
func (c *AppConfig) IsDevelopment() bool {
	return c.Env == "development"
}

// IsProduction 判断是否为生产环境
func (c *AppConfig) IsProduction() bool {
	return c.Env == "production"
}

// IsTest 判断是否为测试环境
func (c *AppConfig) IsTest() bool {
	return c.Env == "test"
}

// GetServerAddr 获取服务器监听地址
func (c *AppConfig) GetServerAddr() string {
	// 如果Host为空，默认使用 "0.0.0.0"
	host := c.Host
	if host == "" {
		host = "0.0.0.0"
	}
	return fmt.Sprintf("%s:%d", host, c.Port)
}
