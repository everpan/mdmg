package main

import (
	"fmt"
	handler2 "github.com/everpan/mdmg/pkg/handler"
	"github.com/everpan/mdmg/pkg/log"
	"github.com/everpan/mdmg/pkg/tenant"
	"github.com/everpan/mdmg/web/handler"
	"github.com/gofiber/contrib/fgprof"
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"time"
	"xorm.io/xorm"
)

func init() {
	viper.SetConfigName("icode")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.icode/")
	viper.SetConfigType("yaml")

	viperDefault()
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
}

func viperDefault() {
	viper.SetDefault("host.addr", "")
	viper.SetDefault("host.port", 8080)
	viper.SetDefault("static.public", "./")
	viper.SetDefault("swagger.file", "./docs/swagger.json")
	viper.SetDefault("swagger.path", "./swagger")
	viper.SetDefault("db.driver", "sqlite3")
	viper.SetDefault("db.connect", "./ic.db")
}
func CreateApp() *fiber.App {

	var defaultEngin, _ = xorm.NewEngine(
		viper.GetString("db.driver"),
		viper.GetString("db.connect"))
	tenant.SetEngine(defaultEngin)

	app := fiber.New()
	logger := log.GetLogger()

	app.Use(fiberzap.New(fiberzap.Config{
		Logger: logger,
	}))

	swgCfg := swagger.Config{
		FilePath: viper.GetString("swagger.file"),
		Path:     viper.GetString("swagger.path"),
	}
	app.Use(swagger.New(swgCfg))

	app.Use(fgprof.New())

	apiRouter := app.Group("/api")
	handler2.AppRouterAdd(apiRouter, &handler.ICoderHandler)

	staticConf := fiber.Static{
		Compress:      true,
		ByteRange:     true,
		Browse:        true,
		Index:         "index.html",
		CacheDuration: 10 * time.Second,
		MaxAge:        3600,
	}

	app.Static(viper.GetString("static.path"),
		viper.GetString("static.root"), staticConf)
	return app
}

func main() {
	app := CreateApp()
	addr := fmt.Sprintf("%s:%d", viper.GetString("host.addr"), viper.GetInt("host.port"))
	app.Listen(addr)
}
