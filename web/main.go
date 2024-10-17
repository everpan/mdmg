package main

import (
	"fmt"
	"github.com/everpan/mdmg/pkg/base/log"
	"github.com/everpan/mdmg/pkg/base/tenant"
	"github.com/everpan/mdmg/pkg/ctx"
	"github.com/everpan/mdmg/pkg/handler"
	"github.com/gofiber/contrib/fgprof"
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"time"
	"xorm.io/xorm"
)

var (
	serverAddress = pflag.String("server.host.address", ":8080",
		"The server address in the format of host:port")
	staticPath = pflag.String("server.public.path", "/", "The path to access public asserts")
	staticRoot = pflag.String("server.public.root", "./static", "The root path of web public")
	dbDriver   = pflag.String("db.driver", "sqlite3", "The database driver name")
	dbConnStr  = pflag.String("db.connect", "./ic_test.db", "The database connection string")
)

func init() {
	viper.BindPFlags(pflag.CommandLine)
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
	viper.SetDefault("swagger.file", "./docs/swagger.json")
	viper.SetDefault("swagger.path", "./swagger")
}
func CreateApp() *fiber.App {

	var defaultEngin, err = xorm.NewEngine(*dbDriver, *dbConnStr)
	if err != nil {
		panic(err)
	}
	tenant.SetSysEngine(defaultEngin)

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
	ctx.AppRouterAdd(apiRouter, &handler.ICoderHandler)
	ctx.AppRouterAddGroup(app, handler.EntityGroupHandler)

	staticConf := fiber.Static{
		Compress:      true,
		ByteRange:     true,
		Browse:        true,
		Index:         "index.html",
		CacheDuration: 10 * time.Second,
		MaxAge:        3600,
	}

	app.Static(*staticPath, *staticRoot, staticConf)
	return app
}

func main() {
	app := CreateApp()
	app.Listen(*serverAddress)
}
