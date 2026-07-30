package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application.
type Config struct {
	AppName              string
	AppEnv               string
	AppPort              string
	AppHost              string
	DBHost               string
	DBPort               string
	DBUser               string
	DBPassword           string
	DBName               string
	DBSSLMode            string
	DBSSLRootCert        string // optional path to a PEM-encoded root CA for verify-full DB connections
	JWTSecret            string
	JWTExpiration        int
	AllowedOrigins       []string
	StripeSecretKey      string
	StripePublishableKey string
	Storage              StorageConfig
	RateReq              float64
	RateBurst            int
}

// StorageConfig selects and configures the file storage provider.
//
// Two providers are supported:
//
//   - "local" (default): writes files to a directory on disk and serves them
//     via the Go server's /uploads/* route. Matches the original behavior.
//
//   - "s3": streams file bytes into any S3-protocol-compatible bucket
//     (AWS S3, Cloudflare R2, MinIO, Supabase Storage, Backblaze B2, ...).
//     Returned image URLs are absolute object URLs served directly by the
//     bucket/CDN, so the Go server stops serving /uploads/* in this mode.
type StorageConfig struct {
	Driver           string // "local" | "s3"
	UploadsDir       string // local mode: directory under cwd
	PublicBaseURL    string // local mode: optional override of /uploads prefix
	S3Bucket         string // s3 mode: bucket name
	S3Region         string // s3 mode: region ("auto" works for R2/MinIO)
	S3Endpoint       string // s3 mode: custom endpoint URL (R2, MinIO, etc.)
	S3AccessKey      string // s3 mode: access key id
	S3SecretKey      string // s3 mode: secret access key
	S3ForcePathStyle bool   // s3 mode: use path-style addressing (MinIO)
	S3PublicBaseURL  string // s3 mode: public origin for returned URLs (CDN/bucket)
}

// Load reads environment variables and returns a Config.
func Load() (*Config, error) {
	_ = godotenv.Load()

	jwtExp, err := strconv.Atoi(getEnv("JWT_EXPIRATION_HOURS", "24"))
	if err != nil {
		jwtExp = 24
	}

	rateReq, err := strconv.Atoi(getEnv("RATE_LIMIT_REQUESTS", "3"))
	if err != nil {
		rateReq = 3
	}
	rateBurst, err := strconv.Atoi(getEnv("RATE_LIMIT_BURST", "5"))
	if err != nil {
		rateBurst = 5
	}

	allowedOrigins := getEnv("ALLOWED_ORIGINS", "http://localhost:5173")

	return &Config{
		AppName:              getEnv("APP_NAME", "Wander"),
		AppEnv:               getEnv("APP_ENV", "development"),
		AppPort:              getEnv("APP_PORT", "8080"),
		AppHost:              getEnv("APP_HOST", "0.0.0.0"),
		DBHost:               getEnv("DB_HOST", "localhost"),
		DBPort:               getEnv("DB_PORT", "5432"),
		DBUser:               getEnv("DB_USER", "wander_user"),
		DBPassword:           getEnv("DB_PASSWORD", "wander_pass"),
		DBName:               getEnv("DB_NAME", "wander_db"),
		DBSSLMode:            getEnv("DB_SSLMODE", "disable"),
		DBSSLRootCert:        getEnv("DB_SSLROOTCERT", ""),
		JWTSecret:            getEnv("JWT_SECRET", "default-secret"),
		JWTExpiration:        jwtExp,
		AllowedOrigins:       strings.Split(allowedOrigins, ","),
		StripeSecretKey:      getEnv("STRIPE_SECRET_KEY", ""),
		StripePublishableKey: getEnv("STRIPE_PUBLISHABLE_KEY", ""),
		RateReq:              float64(rateReq),
		RateBurst:            rateBurst,
		Storage: StorageConfig{
			Driver:           getEnv("STORAGE_DRIVER", "local"),
			UploadsDir:       getEnv("STORAGE_LOCAL_DIR", "uploads"),
			PublicBaseURL:    getEnv("STORAGE_PUBLIC_BASE_URL", "/uploads"),
			S3Bucket:         getEnv("S3_BUCKET", ""),
			S3Region:         getEnv("S3_REGION", "auto"),
			S3Endpoint:       getEnv("S3_ENDPOINT", ""),
			S3AccessKey:      getEnv("S3_ACCESS_KEY", ""),
			S3SecretKey:      getEnv("S3_SECRET_KEY", ""),
			S3ForcePathStyle: getEnvBool("S3_FORCE_PATH_STYLE", false),
			S3PublicBaseURL:  getEnv("S3_PUBLIC_BASE_URL", ""),
		},
	}, nil
}

func (c *Config) DSN() string {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode,
	)
	if c.DBSSLRootCert != "" {
		dsn += fmt.Sprintf(" sslrootcert=%s", c.DBSSLRootCert)
	}
	return dsn
}

// DatabaseURL returns a pgx-compatible connection string.
func (c *Config) DatabaseURL() string {
	u := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		url.QueryEscape(c.DBUser), url.QueryEscape(c.DBPassword),
		c.DBHost, c.DBPort, c.DBName, c.DBSSLMode,
	)
	if c.DBSSLRootCert != "" {
		u += fmt.Sprintf("&sslrootcert=%s", url.QueryEscape(c.DBSSLRootCert))
	}
	return u
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if v := os.Getenv(key); v != "" {
		switch strings.ToLower(v) {
		case "1", "true", "yes", "on":
			return true
		case "0", "false", "no", "off":
			return false
		}
	}
	return fallback
}
