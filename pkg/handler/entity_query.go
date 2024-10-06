package handler

import (
	"fmt"
	"github.com/everpan/mdmg/pkg/ctx"
	"github.com/gofiber/fiber/v2"
)

var EntityQueryHandler = &ctx.IcGroupPathHandler{
	GroupPath: "/entity",
	Handlers: []*ctx.IcPathHandler{
		{
			Path:    "/:className",
			Method:  fiber.MethodPost,
			Handler: query,
		},
	},
}

/*
query

	{
		entity: "user",
		items:{
			"cluster_1":["a","b"],
			"cluster_2":["e","f"]
		},
		alias:{
			"cluster_1": "a",
			"cluster_2": "b",
		},
		where: [
			{"col":"idx","val":3,"op":"gt"},
			{"col":"name,"val":"%ever%","op":"like","combine":"or"},
			{"where":[
				// 内嵌
			]}
		]
	}
*/
func query(c *ctx.IcContext) error {
	var (
		fc        = c.FiberCtx()
		eCtx      = c.EntityCtx()
		className = fc.Params("className")
	)
	eCls, err := eCtx.GetEntityClassByName(className)
	if err != nil {
		return ctx.SendError(fc, fiber.StatusBadRequest, err)
	}
	if eCls == nil {
		err = fmt.Errorf("entity class %s not found", className)
		return ctx.SendError(fc, fiber.StatusBadRequest, err)
	}
	return nil
}
