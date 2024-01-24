package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type DbConf struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	Driver   string
}

type Config struct {
	DbConf
	ApiConfig
	TokenConfig
}

type ApiConfig struct {
	ApiPort string
}

type TokenConfig struct {
	IssuerName       string `json:"issuer_name"`
	JwtSignatureKey  []byte `json:"jwt_signature_key"`
	JwtSigningMethod *jwt.SigningMethodHMAC
	JwtExpiresTime   time.Duration
}

func (c *Config) InitialConfig() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("failed to load env %v", err)
	}

	// Database config
	c.DbConf = DbConf{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB_NAME"),
		Driver:   os.Getenv("DB_DRIVER"),
	}

	c.ApiConfig = ApiConfig{
		ApiPort: os.Getenv("API_PORT"),
	}

	// Config JWT
	tokenExpire, _ := strconv.Atoi(os.Getenv("TOKEN_EXPIRE"))
	c.TokenConfig = TokenConfig{
		IssuerName:       os.Getenv("ISSUER_NAME"),
		JwtSignatureKey:  []byte(os.Getenv("SIGNATURE_KEY")),
		JwtSigningMethod: jwt.SigningMethodHS256,
		JwtExpiresTime:   time.Duration(tokenExpire) * time.Hour,
	}

	if c.Host == "" || c.Port == "" || c.User == "" || c.Password == "" || c.Database == "" || c.Driver == "" || c.IssuerName == "" || c.JwtExpiresTime < 0 || len(c.JwtSignatureKey) == 0 {
		return fmt.Errorf("missing required env variables")
	}
	return nil
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := cfg.InitialConfig(); err != nil {
		return nil, err
	}
	return cfg, nil
}
