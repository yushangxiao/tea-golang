package utils

import (
	"encoding/json"
	"os"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
	jwt.StandardClaims
}

type Config struct {
	AdminSerct []string `json:"adminSerct"`
	JwtKey     string   `json:"jwtKey"`
	RootSerct  string   `json:"rootSerct"`
	Port       string   `json:"port"`
}

var MyConfig *Config

func init() {
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&MyConfig)
	MyConfig.JwtKey = os.Getenv("JWT_KEY")
	MyConfig.RootSerct = os.Getenv("ROOT_SERCT")
	if err != nil {
		panic("failed to read config")
	}
}
