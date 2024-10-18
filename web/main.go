package main

import (
	"fmt"
	"github.com/everpan/mdmg/pkg/base/log"
	"github.com/everpan/mdmg/pkg/config"
	"github.com/everpan/mdmg/pkg/ctx"
	"github.com/everpan/mdmg/pkg/handler"
	"github.com/gofiber/contrib/fgprof"
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"time"
)

var (
	logger        = log.GetLogger()
	serverAddress = pflag.String("server.host-address", ":8080",
		"The server address in the format of host:port")
	staticPath = pflag.String("server.public-path", "/", "The path to access public asserts")
	staticRoot = pflag.String("server.public-root", "./static", "The root path of web public")
	//dbDriver    = pflag.String("server.db-driver", "sqlite3", "The database driver name")
	//dbConnStr   = pflag.String("server.db-connect", "./ic_test.db", "The database connection string")
	swaggerFile = pflag.String("server.swagger-file", "./docs/swagger.json", "The swagger file to serve")
	swaggerPath = pflag.String("server.swagger-path", "/swagger/", "The swagger url path to serve")
)

func updateConfigValues() error {
	*serverAddress = viper.GetString("server.host-address")
	*staticRoot = viper.GetString("server.public-root")
	*staticPath = viper.GetString("server.public-path")
	//*dbDriver = viper.GetString("server.db-driver")
	//*dbConnStr = viper.GetString("server.db-connect")
	*swaggerFile = viper.GetString("server.swagger-file")
	*swaggerPath = viper.GetString("server.swagger-path")
	return nil
}

func init() {
	viper.BindPFlags(pflag.CommandLine)
	viper.SetConfigName("icode")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.icode/")
	viper.SetConfigType("yaml")
	// create yaml config
	// viper.WriteConfigAs("icode.yaml")
	viper.SafeWriteConfigAs("icode.yaml")

	config.RegisterReloadViperFunc(updateConfigValues)
	err := viper.ReadInConfig()

	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
	logger.Info("read in config finished",
		zap.String("server.swagger-file", viper.GetString("server.swagger-file")))
}

func CreateApp() *fiber.App {
	logger.Debug("variables", zap.String("serverAddress", *serverAddress),
		zap.String("staticPath", *staticPath), zap.String("staticRoot", *staticRoot),
		// zap.String("dbDriver", *dbDriver), zap.String("dbConnStr", *dbConnStr),
		zap.String("swaggerFile", *swaggerFile), zap.String("swaggerPath", *swaggerPath))
	// config.AcquireEngine(*dbDriver, *dbConnStr)

	app := fiber.New()

	app.Use(fiberzap.New(fiberzap.Config{Logger: logger}))

	swgCfg := swagger.Config{FilePath: *swaggerFile, Path: *swaggerPath}
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
	pflag.Parse()
	config.ReloadViperConfig()
	app := CreateApp()
	app.Listen(*serverAddress)
}
