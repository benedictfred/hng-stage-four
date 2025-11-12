package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	RabbitMQ RabbitMQConfig `mapstructure:"rabbitmq"`
	FCM      FCMConfig      `mapstructure:"fcm"`
	Log      LogConfig      `mapstructure:"log"`
	Queue    QueueConfig    `mapstructure:"queue"`
}

type ServerConfig struct {
	Port            string        `mapstructure:"port"`
	Mode            string        `mapstructure:"mode"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

type DatabaseConfig struct {
	Host            string        `mapstructure:"host"`
	Port            string        `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	Name            string        `mapstructure:"name"`
	SSLMode         string        `mapstructure:"ssl_mode"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type FCMConfig struct {
	CredentialsJSON string `mapstructure:"credentials_json"`
	ProjectID       string `mapstructure:"project_id"`
	UseFile         bool   `mapstructure:"use_file"`
}

type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

type RabbitMQConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	VHost    string `mapstructure:"vhost"`
}

type QueueConfig struct {
	Worker      WorkerConfig      `mapstructure:"worker"`
	Retry       RetryConfig       `mapstructure:"retry"`
	Validation  ValidationConfig  `mapstructure:"validation"`
}

type WorkerConfig struct {
	PrefetchCount int           `mapstructure:"prefetch_count"`
	PollInterval  time.Duration `mapstructure:"poll_interval"`
	BatchSize     int           `mapstructure:"batch_size"`
}

type RetryConfig struct {
	MaxRetries int           `mapstructure:"max_retries"`
	Backoff    time.Duration `mapstructure:"backoff"`
}

type ValidationConfig struct {
	Enabled bool          `mapstructure:"enabled"`
	Timeout time.Duration `mapstructure:"timeout"`
}

func Load() (*Config, error) {
	// Load .env file if exists
	godotenv.Load() // This will load .env file, but doesn't fail if it doesn't exist

	// You can also load specific environment files
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}
	godotenv.Load(fmt.Sprintf(".env.%s", env))

	// Set up Viper
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/etc/push-service/")

	// Set defaults
	setDefaults()

	// Read config file (optional - can use env vars only)
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Warning: Config file not found, using environment variables: %v\n", err)
	}

	// Bind environment variables
	bindEnvVars()

	// Unmarshal config
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	// Validate required fields
	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func setDefaults() {
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("server.shutdown_timeout", "30s")

	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", "5432")
	viper.SetDefault("database.name", "push_service")
	viper.SetDefault("database.ssl_mode", "disable")
	viper.SetDefault("database.max_open_conns", 25)
	viper.SetDefault("database.max_idle_conns", 25)
	viper.SetDefault("database.conn_max_lifetime", "5m")

	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", "6379")
	viper.SetDefault("redis.db", 0)

	viper.SetDefault("rabbitmq.host", "localhost")
	viper.SetDefault("rabbitmq.port", "5672")
	viper.SetDefault("rabbitmq.username", "guest")
	viper.SetDefault("rabbitmq.password", "guest")
	viper.SetDefault("rabbitmq.vhost", "/")

	viper.SetDefault("queue.worker.prefetch_count", 10)
	viper.SetDefault("queue.worker.poll_interval", "1s")
	viper.SetDefault("queue.worker.batch_size", 10)
	viper.SetDefault("queue.retry.max_retries", 5)
	viper.SetDefault("queue.retry.backoff", "5s")
	viper.SetDefault("queue.validation.enabled", true)
	viper.SetDefault("queue.validation.timeout", "5s")

	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "json")
}

func bindEnvVars() {
	// Server
	viper.BindEnv("server.port", "SERVER_PORT")
	viper.BindEnv("server.mode", "SERVER_MODE")
	viper.BindEnv("server.shutdown_timeout", "SERVER_SHUTDOWN_TIMEOUT")

	// Database
	viper.BindEnv("database.host", "DB_HOST")
	viper.BindEnv("database.port", "DB_PORT")
	viper.BindEnv("database.user", "DB_USER")
	viper.BindEnv("database.password", "DB_PASSWORD")
	viper.BindEnv("database.name", "DB_NAME")
	viper.BindEnv("database.ssl_mode", "DB_SSL_MODE")
	viper.BindEnv("database.max_open_conns", "DB_MAX_OPEN_CONNS")
	viper.BindEnv("database.max_idle_conns", "DB_MAX_IDLE_CONNS")
	viper.BindEnv("database.conn_max_lifetime", "DB_CONN_MAX_LIFETIME")

	// Redis
	viper.BindEnv("redis.host", "REDIS_HOST")
	viper.BindEnv("redis.port", "REDIS_PORT")
	viper.BindEnv("redis.password", "REDIS_PASSWORD")
	viper.BindEnv("redis.db", "REDIS_DB")

	// RabbitMQ
	viper.BindEnv("rabbitmq.host", "RABBITMQ_HOST")
	viper.BindEnv("rabbitmq.port", "RABBITMQ_PORT")
	viper.BindEnv("rabbitmq.username", "RABBITMQ_USERNAME")
	viper.BindEnv("rabbitmq.password", "RABBITMQ_PASSWORD")
	viper.BindEnv("rabbitmq.vhost", "RABBITMQ_VHOST")

	// Queue
	viper.BindEnv("queue.worker.prefetch_count", "QUEUE_WORKER_PREFETCH_COUNT")
	viper.BindEnv("queue.worker.poll_interval", "QUEUE_WORKER_POLL_INTERVAL")
	viper.BindEnv("queue.worker.batch_size", "QUEUE_WORKER_BATCH_SIZE")
	viper.BindEnv("queue.retry.max_retries", "QUEUE_RETRY_MAX_RETRIES")
	viper.BindEnv("queue.retry.backoff", "QUEUE_RETRY_BACKOFF")
	viper.BindEnv("queue.validation.enabled", "QUEUE_VALIDATION_ENABLED")
	viper.BindEnv("queue.validation.timeout", "QUEUE_VALIDATION_TIMEOUT")

	// FCM
	viper.BindEnv("fcm.credentials_json", "FCM_CREDENTIALS_JSON")
	viper.BindEnv("fcm.project_id", "FCM_PROJECT_ID")
	viper.BindEnv("fcm.use_file", "FCM_USE_FILE")

	// Log
	viper.BindEnv("log.level", "LOG_LEVEL")
	viper.BindEnv("log.format", "LOG_FORMAT")
}

// GetDatabaseURL builds the database connection URL
func (db *DatabaseConfig) GetDatabaseURL() string {
	if url := os.Getenv("DATABASE_URL"); url != "" {
		return url
	}

	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		db.User,
		db.Password,
		db.Host,
		db.Port,
		db.Name,
		db.SSLMode,
	)
}

// GetRedisURL builds the Redis connection URL
func (redis *RedisConfig) GetRedisURL() string {
	if url := os.Getenv("REDIS_URL"); url != "" {
		return url
	}

	if redis.Password != "" {
		return fmt.Sprintf("redis://:%s@%s:%s/%d",
			redis.Password,
			redis.Host,
			redis.Port,
			redis.DB,
		)
	}

	return fmt.Sprintf("redis://%s:%s/%d",
		redis.Host,
		redis.Port,
		redis.DB,
	)
}

func validateConfig(config *Config) error {
	// Validate required fields
	if config.Database.User == "" {
		return fmt.Errorf("database user is required")
	}
	if config.Database.Password == "" {
		return fmt.Errorf("database password is required")
	}
	if config.FCM.CredentialsJSON == "" {
		return fmt.Errorf("FCM credentials are required")
	}

	return nil
}

// GetFCMCredentials returns FCM credentials as byte array
func (c *FCMConfig) GetFCMCredentials() ([]byte, error) {
	if c.UseFile {
		// Read from file
		return os.ReadFile(c.CredentialsJSON)
	} else {
		// Treat as JSON string
		return []byte(c.CredentialsJSON), nil
	}
}
