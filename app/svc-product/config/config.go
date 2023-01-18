package config

import (
	"log"
	"strings"

	"github.com/dewidyabagus/monorepo/datasource/sql"
	"github.com/dewidyabagus/monorepo/shared-libs/envar"
	"github.com/joho/godotenv"
)

type Group struct {
	Env      string
	AppName  string
	AppPort  int
	Database sql.MySQLConfig
}

func LoadConfig(filenames ...string) *Group {
	if len(filenames) == 0 {
		filenames = []string{".env"}
	}

	var err error
	for _, file := range filenames {
		if err = godotenv.Load(file); err == nil {
			break
		}
	}
	if err != nil {
		log.Printf("%v not found \n", strings.Join(filenames, " "))
	}

	env := envar.GetEnv("ENV", "dev")
	return &Group{
		Env:     env,
		AppName: envar.GetEnv("APP_NAME", "svc-product"),
		AppPort: envar.GetEnv("APP_PORT", 8002),
		Database: sql.MySQLConfig{
			Host:             envar.GetEnv("DATABASE_HOST", "localhost"),
			Port:             envar.GetEnv("DATABASE_PORT", 3306),
			User:             envar.GetEnv("DATABASE_USERNAME", "root"),
			Password:         envar.GetEnv("DATABASE_PASSWORD", ""),
			Schema:           envar.GetEnv("DATABASE_SCHEMA", "db"),
			MaxIdleConns:     envar.GetEnv("DATABASE_MAX_IDLE", 5),
			MaxOpenConns:     envar.GetEnv("DATABASE_MAX_CONN", 5),
			ConnMaxLifetime:  envar.GetEnv("DATABASE_CONN_LIFETIME", 180),
			Environment:      env,
			SlowLogThreshold: envar.GetEnv("DATABASE_SLOW_LOG_THRESHOLD", 200),
		},
	}
}
