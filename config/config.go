package config

import (
	"fmt"
	"go/build"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

var ServiceName = "chat_app_server"

type Secrets struct {
	Db_User       string `json:"USER" envconfig:"DB_USER"`
	DbName        string `json:"DBNAME" envconfig:"DBNAME"`
	Password      string `json:"PASSWORD" envconfig:"PASSWORD"`
	Host          string `json:"HOST" envconfig:"HOST"`
	Port          string `json:"PORT" envconfig:"PORT"`
	JwtSecret     string `json:"JWT_SECRET" envconfig:"JWT_SECRET"`
	Neo4jUri      string `json:"NEO4J_URI" envconfig:"NEO4J_URI"`
	Neo4jUser     string `json:"NEO4J_USER" envconfig:"NEO4J_USER"`
	Neo4jPassword string `json:"NEO4J_PASSWORD" envconfig:"NEO4J_PASSWORD"`
	RedisURL      string `json:"REDIS_URL" envconfig:"REDIS_URL"`
}

var ss Secrets

func init() {
	importPath := fmt.Sprintf("%s/config", ServiceName)
	p, err := build.Default.Import(importPath, "", build.FindOnly)
	if err == nil {
		env := filepath.Join(p.Dir, "../.env")
		_ = godotenv.Load(env)
	}
	ss = Secrets{}
	ss.DbName = os.Getenv("DBNAME")
	ss.Db_User = os.Getenv("DB_USER")
	ss.Password = os.Getenv("PASSWORD")
	ss.Host = os.Getenv("HOST")
	ss.Port = os.Getenv("PORT")
	ss.JwtSecret = os.Getenv("JWT_SECRET")
	ss.Neo4jUri = strings.TrimSpace(os.Getenv("NEO4J_URI"))
	ss.Neo4jUser = os.Getenv("NEO4J_USER")
	ss.Neo4jPassword = os.Getenv("NEO4J_PASSWORD")
	ss.RedisURL = strings.TrimSpace(os.Getenv("REDIS_URL"))
}

func GetSecrets() Secrets {
	return ss
}
