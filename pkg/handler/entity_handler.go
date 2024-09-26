package handler

import (
	"fmt"
	"github.com/everpan/mdmg/pkg/base/entity"
	"github.com/everpan/mdmg/pkg/ctx"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"strings"
)

var EntityGroupHandler = &ctx.IcGroupPathHandler{
	GroupPath: "/entity",
	Handlers: []*ctx.IcPathHandler{
		{
			Path:    "/meta/list/:page?", // page like 2-20, pageNum = 2, pageSize: 20
			Method:  fiber.MethodGet,
			Handler: metaList,
		},
		{
			Path:    "/meta/class/:classNameOrId?",
			Method:  fiber.MethodGet,
			Handler: metaDetail,
		},
	},
}

func metaDetail(c *ctx.IcContext) error {
	var (
		meta          = entity.IcEntityMeta{}
		err           error
		fc            = c.FiberCtx()
		classNameOrId = fc.Params("classNameOrId")
	)

	if classNameOrId == "" {
		return ctx.SendError(fc, fiber.StatusBadRequest,
			fmt.Errorf("class name or id not specified"))
	}
	classId, err := strconv.ParseUint(classNameOrId, 10, 32)
	if classId == 0 && err == nil {
		return ctx.SendError(fc, fiber.StatusBadRequest,
			fmt.Errorf("classId=%d is required and must be gt zero", classId))
	}
	entityCtx := c.EntityCtx()
	if err != nil { // classNameOrId name
		meta.EntityClass, err = entityCtx.GetEntityClassByName(classNameOrId)
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
	pageInfo := fc.Params("page")
	if pageInfo == "" {
		c.Page.Reset()
	} else {
		sp := strings.Split(pageInfo, "-")
		c.Page.No, _ = strconv.Atoi(sp[0])
		c.Page.Size = 20
		if len(sp) > 1 {
			c.Page.Size, _ = strconv.Atoi(sp[1])
		}
	}
	// fmt.Printf("-- page %v\n", c.Page)
	count, _ := c.Engine().Count(&entity.IcEntityClass{TenantId: c.Tenant().Idx})
	offset := c.Page.CalCountOffset(int(count))

	var eClasses []*entity.IcEntityClass
	e := c.Engine().Limit(c.Page.Size, offset).Where("tenant_id = ?", c.Tenant().Idx).Find(&eClasses)
	if nil != e {
		return ctx.SendError(fc, fiber.StatusInternalServerError, e)
	}
	metas := make([]*entity.IcEntityMeta, len(eClasses))
	for i, class := range eClasses {
		tables, _ := c.EntityCtx().GetClusterTables(class.ClassId)
		metas[i] = &entity.IcEntityMeta{
			EntityClass:   class,
			ClusterTables: tables,
		}
	}
	/*	// using where in
		in := make([]uint32, len(eClasses))
		for i, class := range eClasses {
			in[i] = class.ClassId
		}
		clusterTables := make([]*entity.IcClusterTable, 0)
		e = c.Engine().Table("ic_cluster_table").In("class_id", in).Find(clusterTables)
		if nil != e {
			return ctx.SendError(fc, fiber.StatusInternalServerError, e)
		}
	*/
	return ctx.SendSuccessWithPage(fc, metas, *c.Page)
}
