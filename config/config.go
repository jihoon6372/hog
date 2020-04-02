package config

import (
	"github.com/dgrijalva/jwt-go"
)

// JwtCustomClaims are custom claims extending default ones.
type JwtCustomClaims struct {
	Username string
	Email    string
	jwt.StandardClaims
}

// Config ...
type Config struct {
	Database struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Username string `yaml:"user"`
		Password string `yaml:"password"`
		Database string `yarm:"database"`
		SSLMode  string `yarm:"sslmode"`
	} `yaml:"database"`
	Config struct {
		SecretKey string `yaml:"secretkey"`
	} `yaml:"config"`
}
