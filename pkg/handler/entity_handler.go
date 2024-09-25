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
			Path:    "/meta/list",
			Method:  fiber.MethodGet,
			Handler: metaList,
		},
		{
			Path:    "/meta/:class?",
			Method:  fiber.MethodGet,
			Handler: metaDetail,
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
	fc := c.FiberCtx()
	// sql := `select a.*,b.* from ic_entity_class as a, ic_cluster_table as b where a.class_id = b.class_id`
	offset := c.Page.Number * c.Page.Size
	// r, e := c.Engine().Limit(c.Page.Size, offset).SQL(sql).QueryInterface()
	// 以上方式 limit不起效果
	var eClasses []*entity.IcEntityClass
	e := c.Engine().Limit(c.Page.Size, offset).Find(&eClasses)
	if nil != e {
		return ctx.SendError(fc, fiber.StatusInternalServerError, e)
	}
	// todo sql using in (....) ?
	metas := make([]*entity.IcEntityMeta, len(eClasses))
	for i, class := range eClasses {
		tables, _ := c.EntityCtx.GetClusterTables(class.ClassId)
		metas[i] = &entity.IcEntityMeta{
			EntityClass:   class,
			ClusterTables: tables,
		}
	}
	return ctx.SendSuccess(fc, metas)
}
