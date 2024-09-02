package main

import (
	"fmt"
	"github.com/everpan/mdmg/v8runtime"
	"github.com/everpan/mdmg/web/handler"
	"github.com/gofiber/contrib/fgprof"
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"time"
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
}
func CreateApp() *fiber.App {
	app := fiber.New()
	// contrib/fiberzap
	logger, _ := zap.NewDevelopment()
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
	apiRouter.Group(handler.ICoderHandler.Path, handler.ICoderHandler.Handler)

	staticConf := fiber.Static{
		Compress:      true,
		ByteRange:     true,
		Browse:        true,
		Index:         "index.html",
		CacheDuration: 10 * time.Second,
		MaxAge:        3600,
	}

	app.Static("/", "./public", staticConf)
	return app
}

func main() {
	defer v8runtime.DisposeCtxPool()

	app := CreateApp()
	addr := fmt.Sprintf("%s:%d", viper.GetString("host.addr"), viper.GetInt("host.port"))
	app.Listen(addr)
}
