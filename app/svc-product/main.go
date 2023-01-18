package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dewidyabagus/monorepo/app/svc-product/config"
	"github.com/dewidyabagus/monorepo/datasource/sql"
	"github.com/dewidyabagus/monorepo/shared-libs/slack"
	"github.com/labstack/echo/v4"
)

var configs = []string{
	"./config/.env",
	"./app/svc-product/config/.env",
}

func main() {
	cfg := config.LoadConfig(configs...)

	db, err := sql.NewMySQL(&cfg.Database)
	if err != nil {
		log.Fatalln("new session:", err.Error())
	}
	inst, _ := db.DB()
	defer inst.Close()

	// Example untuk penerapan package shared-libs
	_ = slack.NewSlackApi(&slack.SlackOption{})

	e := echo.New()

	e.GET("/ping", func(ctx echo.Context) error {
		if err := inst.PingContext(ctx.Request().Context()); err != nil {
			return ctx.JSON(http.StatusInternalServerError, echo.Map{"message": err.Error()})
		}
		return ctx.JSON(http.StatusOK, echo.Map{"message": "OK"})
	})

	log.Fatal(e.Start(fmt.Sprintf("0.0.0.0:%d", cfg.AppPort)))
}
