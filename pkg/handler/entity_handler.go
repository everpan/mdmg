package handler

import (
	"fmt"
	"github.com/everpan/mdmg/pkg/base/entity"
	"github.com/everpan/mdmg/pkg/ctx"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

var EntityGroupHandler = &ctx.IcGroupPathHandler{
	GroupPath: "/entity",
	Handlers: []*ctx.IcPathHandler{
		{
			Path:    "/meta/:class?",
			Method:  fiber.MethodGet,
			Handler: metaDetail,
		},
		{
			Path:    "/meta/list",
			Handler: metaList,
		},
	},
}

func metaDetail(c *ctx.IcContext) error {
	var (
		meta  = entity.IcEntityMeta{}
		err   error
		fc    = c.FiberCtx()
		class = fc.Params("class")
	)

	if class == "" {
		return ctx.SendError(fc, fiber.StatusBadRequest,
			fmt.Errorf("class not specified"))
	}
	classId, err := strconv.ParseUint(class, 10, 32)
	if classId == 0 && err == nil {
		return ctx.SendError(fc, fiber.StatusBadRequest,
			fmt.Errorf("classId=%d is required and must be gt zero", classId))
	}
	entityCtx := c.EntityCtx
	if err != nil { // class name
		meta.EntityClass, err = entityCtx.GetEntityClassByName(class)
	} else {
		meta.EntityClass, err = entityCtx.GetEntityClassById(uint32(classId))
	}
	if nil != err {
		return ctx.SendError(fc, fiber.StatusBadRequest, err)
	}
	meta.ClusterTables, err = entityCtx.GetClusterTables(meta.EntityClass.ClassId)
	if nil != err {
		return ctx.SendError(fc, fiber.StatusBadRequest, err)
	}
	return ctx.SendSuccess(fc, meta)
}
func metaList(c *ctx.IcContext) error {
	//
	return nil
}
