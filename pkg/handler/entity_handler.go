package handler

import (
	"fmt"
	"github.com/everpan/mdmg/pkg/ctx"
	"github.com/gofiber/fiber/v2"
)

var EntityGroupHandler = &ctx.IcGroupPathHandler{
	GroupPath: "/entity",
	Handlers: []*ctx.IcPathHandler{
		{
			Path:    "/meta/:classId",
			Handler: metaDetail,
		},
	},
}

func metaDetail(c *ctx.IcContext) error {
	fb := c.FiberCtx()
	classId, _ := fb.ParamsInt("classId", 0)
	if classId == 0 {
		return ctx.SendError(fb, fiber.StatusBadRequest, fmt.Errorf("classId=%d is required and must be gt zero", classId))
	}
	// entity.
	return nil
}
