package config

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	OSS      OSSConfig      `mapstructure:"oss"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Log      LogConfig      `mapstructure:"log"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`       // debug, info, warn, error
	Format     string `mapstructure:"format"`      // json, text
	Output     string `mapstructure:"output"`      // stdout, file
	FilePath   string `mapstructure:"file_path"`   // log file path
	MaxSize    int    `mapstructure:"max_size"`    // max size in MB
	MaxBackups int    `mapstructure:"max_backups"` // max backup files
	MaxAge     int    `mapstructure:"max_age"`     // max age in days
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	Charset  string `mapstructure:"charset"`
}

type JWTConfig struct {
	Secret string `mapstructure:"secret"`
	Expire int    `mapstructure:"expire"`
}

type OSSConfig struct {
	Endpoint        string `mapstructure:"endpoint"`
	AccessKeyID     string `mapstructure:"access_key_id"`
	AccessKeySecret string `mapstructure:"access_key_secret"`
	BucketName      string `mapstructure:"bucket_name"`
	BaseURL         string `mapstructure:"base_url"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

var AppConfig *Config
var DB *gorm.DB

func LoadConfig(configPath string) error {
	if configPath != "" {
		dir := filepath.Dir(configPath)
		filename := filepath.Base(configPath)
		ext := filepath.Ext(filename)
		name := strings.TrimSuffix(filename, ext)

		viper.SetConfigName(name)
		viper.SetConfigType(strings.TrimPrefix(ext, "."))
		viper.AddConfigPath(dir)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
	}

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	AppConfig = &Config{}
	if err := viper.Unmarshal(AppConfig); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Set defaults for log config
	if AppConfig.Log.Level == "" {
		AppConfig.Log.Level = "info"
	}
	if AppConfig.Log.Format == "" {
		AppConfig.Log.Format = "text"
	}
	if AppConfig.Log.Output == "" {
		AppConfig.Log.Output = "stdout"
	}

	return nil
}

func InitDatabase() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		AppConfig.Database.User,
		AppConfig.Database.Password,
		AppConfig.Database.Host,
		AppConfig.Database.Port,
		AppConfig.Database.DBName,
		AppConfig.Database.Charset,
	)

	var logLevel logger.LogLevel
	if AppConfig.Server.Mode == "debug" {
		logLevel = logger.Info
	} else {
		logLevel = logger.Silent
	}

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	log.Println("Database connected successfully")
	return nil
}
