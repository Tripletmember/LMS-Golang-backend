package config

import (
	"os"
	"time"

	"github.com/spf13/viper"
)

const (
	defaultHTTPPort               = "8000"
	defaultHttpRWTimeout          = 10 * time.Second
	defaultHttpMaxHeaderMegabytes = 1
	defaultAccessTokenTTL         = 15 * time.Minute
	defaultRefreshTokenTTL        = 24 * time.Hour * 30
	defaultLimiterRPS             = 10
	defaultLimiterBurst           = 2
	defaultLimiterTTL             = 10 * time.Minute
	defaultVerificationCodeLength = 8

	EnvLocal = "local"
	Prod     = "prod"
)

type (
	Config struct {
		Environment string
		Mongo       MongoConfig
		HTTP        HTTPConfig
		Auth        AuthConfig
		FileStorage FileStorageConfig
		Email       EmailConfig
		Payment     PaymentConfig
		Limiter     LimiterConfig
		CacheTTL    time.Duration `mapstructure:"ttl"`
		FrontendURL string
		SMTP        SMTPConfig
		Cloudflare  CloudflareConfig
	}

	MongoConfig struct {
		URI      string
		User     string
		Password string
		Name     string `mapstructure:"databaseName"`
	}

	AuthConfig struct {
		JWT                    JWTConfig
		PasswordSalt           string
		VerificationCodeLength int `mapstructure:"verificationCodeLength"`
	}

	JWTConfig struct {
		AccessTokenTTL  time.Duration `mapstructure:"accessTokenTTL"`
		RefreshTokenTTL time.Duration `mapstructure:"refreshTokenTTL"`
		SigningKey      string
	}

	FileStorageConfig struct {
		Endpoint  string
		Bucket    string
		AccessKey string
		SecretKey string
	}

	EmailConfig struct {
		SendPulse SendPulseConfig
		Templates EmailTemplates
		Subjects  EmailSubjects
	}

	SendPulseConfig struct {
		ListID       string
		ClientID     string
		ClientSecret string
	}

	EmailTemplates struct {
		Verification       string `mapstructure:"verification_email"`
		PurchaseSuccessful string `mapstructure:"purchase_successful"`
	}

	EmailSubjects struct {
		Verification       string `mapstructure:"verification_email"`
		PurchaseSuccessful string `mapstructure:"purchase_successful"`
	}

	PaymentConfig struct {
		Fondy       FondyConfig
		CallbackURL string
		ResponseURL string
	}

	FondyConfig struct {
		MerchantId       string
		MerchantPassword string
	}

	HTTPConfig struct {
		Host               string        `mapstructure:"host"`
		Port               string        `mapstructure:"port"`
		ReadTimeout        time.Duration `mapstructure:"readTimeout"`
		WriteTimeout       time.Duration `mapstructure:"writeTimeout"`
		MaxHeaderMegabytes int           `mapstructure:"maxHeaderBytes"`
	}

	LimiterConfig struct {
		RPS   int
		Burst int
		TTL   time.Duration
	}

	SMTPConfig struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
		From string `mapstructure:"from"`
		Pass string
	}

	CloudflareConfig struct {
		ApiKey      string
		Email       string
		ZoneEmail   string
		CnameTarget string
	}
)

// Init populates Config struct with values from config file
// located at filepath and environment variables.
func Init(configsDir string) (*Config, error) {
	populateDefaults()

	if err := parseConfigFile(configsDir, os.Getenv("APP_ENV")); err != nil {
		return nil, err
	}

	var cfg Config

	setFromEnv(&cfg)

	return &cfg, nil
}

func setFromEnv(cfg *Config) {
	cfg.Mongo.URI = os.Getenv("MONGO_URI")
	cfg.Mongo.User = os.Getenv("MONGO_USER")
	cfg.Mongo.Password = os.Getenv("MONGO_PASS")

	cfg.Auth.PasswordSalt = os.Getenv("PASSWORD_SALT")
	cfg.Auth.JWT.SigningKey = os.Getenv("JWT_SIGNING_KEY")

	cfg.Email.SendPulse.ClientSecret = os.Getenv("SENDPULSE_SECRET")
	cfg.Email.SendPulse.ClientID = os.Getenv("SENDPULSE_ID")
	cfg.Email.SendPulse.ListID = os.Getenv("SENDPULSE_LISTID")

	cfg.HTTP.Host = os.Getenv("HTTP_HOST")

	cfg.Payment.Fondy.MerchantId = os.Getenv("FONDY_MERCHANT_ID")
	cfg.Payment.Fondy.MerchantPassword = os.Getenv("FONDY_MERCHANT_PASS")
	cfg.Payment.CallbackURL = os.Getenv("PAYMENT_CALLBACK_URL")
	cfg.Payment.ResponseURL = os.Getenv("PAYMENT_REDIRECT_URL")

	cfg.FrontendURL = os.Getenv("FRONTEND_URL")

	cfg.SMTP.Pass = os.Getenv("SMTP_PASSWORD")

	cfg.Environment = os.Getenv("APP_ENV")

	cfg.FileStorage.Endpoint = os.Getenv("STORAGE_ENDPOINT")
	cfg.FileStorage.AccessKey = os.Getenv("STORAGE_ACCESS_KEY")
	cfg.FileStorage.SecretKey = os.Getenv("STORAGE_SECRET_KEY")
	cfg.FileStorage.Bucket = os.Getenv("STORAGE_BUCKET")

	cfg.Cloudflare.ApiKey = os.Getenv("CLOUDFLARE_API_KEY")
	cfg.Cloudflare.Email = os.Getenv("CLOUDFLARE_EMAIL")
	cfg.Cloudflare.ZoneEmail = os.Getenv("CLOUDFLARE_ZONE_EMAIL")
	cfg.Cloudflare.CnameTarget = os.Getenv("CLOUDFLARE_CNAME_TARGET")
}

func parseConfigFile(folder, env string) error {
	viper.AddConfigPath(folder)
	viper.SetConfigName("main")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if env == EnvLocal {
		return nil
	}

	viper.SetConfigName(env)

	return viper.MergeInConfig()
}

func populateDefaults() {
	if _, ok := os.LookupEnv("HTTP_HOST"); ok {
		os.Setenv("HTTP_HOST", defaultHTTPPort)
	}

	viper.SetDefault("http.max_header_megabytes", defaultHttpMaxHeaderMegabytes)
	viper.SetDefault("http.timeouts.read", defaultHttpRWTimeout)
	viper.SetDefault("http.timeouts.write", defaultHttpRWTimeout)
	viper.SetDefault("auth.accessTokenTTL", defaultAccessTokenTTL)
	viper.SetDefault("auth.refreshTokenTTL", defaultRefreshTokenTTL)
	viper.SetDefault("auth.verificationCodeLength", defaultVerificationCodeLength)
	viper.SetDefault("limiter.rps", defaultLimiterRPS)
	viper.SetDefault("limiter.burst", defaultLimiterBurst)
	viper.SetDefault("limiter.ttl", defaultLimiterTTL)
}
