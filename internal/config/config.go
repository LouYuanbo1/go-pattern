package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	Database   DatabaseConfig   `mapstructure:"database"`
	Redis      RedisConfig      `mapstructure:"redis"`
	LocalCache LocalCacheConfig `mapstructure:"local_cache"`
}

type DatabaseConfig struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	User            string `mapstructure:"user"`
	Password        string `mapstructure:"password"`
	DBName          string `mapstructure:"dbname"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`    // 最大打开连接数 (建议值: 25)
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`    // 最大空闲连接数 (建议值: 25)
	ConnMaxLifetime string `mapstructure:"conn_max_lifetime"` // 连接最大生命周期 (建议值: 5m)
	TimeZone        string `mapstructure:"time_zone"`         // 时区配置 (示例: "Asia/Shanghai")
	//Silent=1, Warn=2, Error=3, Info=4 (默认)
	LogLevel int `mapstructure:"log_level"` // 日志级别 (示例: "info")
	//SchemaFile      string `mapstructure:"schema_file"`       // schema.sql 文件路径
}

type RedisConfig struct {
	Host          string `mapstructure:"host"`
	Port          int    `mapstructure:"port"`
	Password      string `mapstructure:"password"`
	DB            int    `mapstructure:"db"`
	Protocol      int    `mapstructure:"protocol"`
	UnstableResp3 bool   `mapstructure:"unstable_resp3"`
	MaxSize       int    `mapstructure:"max_size"`
	DefaultTTL    int64  `mapstructure:"default_ttl"`
}

type LocalCacheConfig struct {
	NumCounters int64 `mapstructure:"num_counters"`
	MaxCost     int64 `mapstructure:"max_cost"`
	BufferItems int64 `mapstructure:"buffer_items"`
	DefaultTTL  int64 `mapstructure:"default_ttl"`
}

type JWTConfig struct {
	SecretKey string `mapstructure:"secret_key"`
	// 令牌过期时间(单位:小时)
	TokenExpire int64 `mapstructure:"token_expire"`
}

func InitConfig() (*Config, error) {
	// 配置文件设置
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	path := filepath.Join("..", "config")
	//Viper查询的路径是相对于当前工作目录（current working directory） 的，而不是相对于可执行文件的位置或源代码文件的位置。
	// 1.尝试加载main.go所在目录的上两级目录下config.yaml
	viper.AddConfigPath(path)
	// 2.尝试加载当前工作目录下的config.yaml
	viper.AddConfigPath(".")

	exePath, _ := os.Executable()
	exeDir := filepath.Dir(exePath)

	// 设置绝对路径
	// 3.尝试加载main.go可执行文件所在目录的上两级目录下的config.yaml
	viper.AddConfigPath(filepath.Join(exeDir, "..", "..", "config"))
	// 4.尝试加载可执行文件所在目录下的config.yaml
	viper.AddConfigPath(exeDir)

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("No config file found, using defaults")
		} else {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	// 绑定环境变量
	viper.AutomaticEnv()

	// 监控配置变化
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
	viper.WatchConfig()

	cfg, err := parseConfig()
	if err != nil {
		log.Fatalf("解析配置文件失败: %v", err)
	}
	return cfg, nil
}

func parseConfig() (*Config, error) {
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %w", err)
	}
	return &cfg, nil
}
