package handler

import (
	"github.com/everpan/mdmg/pkg/config"
	"github.com/everpan/mdmg/pkg/ctx"
	"github.com/gofiber/fiber/v2"
)

var ConfigHandler = &ctx.IcGroupPathHandler{
	GroupPath: "/config",
	Handlers: []*ctx.IcPathHandler{
		{
			Path:    "/config",
			Method:  fiber.MethodGet,
			Handler: getConfig,
		},
		{
			Path:    "/schema",
			Method:  fiber.MethodGet,
			Handler: exportConfig,
		},
	},
}

func getConfig(c *ctx.IcContext) error {
	return nil
}

func exportConfig(c *ctx.IcContext) error {
	schema := config.GlobalConfig.ExportSchema()
	return ctx.SendSuccess(c.FiberCtx(), schema)
}
