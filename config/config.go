package config

import (
	"fmt"
	"go/build"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

var ServiceName = "chat_app_server"

type Secrets struct {
	Db_User  string `json:"USER" envconfig:"DB_USER"`
	DbName   string `json:"DBNAME" envconfig:"DBNAME"`
	Password string `json:"PASSWORD" envconfig:"PASSWORD"`
	Host     string `json:"HOST" envconfig:"HOST"`
	Port     string `json:"PORT" envconfig:"PORT"`
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
}

func GetSecrets() Secrets {
	return ss
}
