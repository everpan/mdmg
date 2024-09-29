package handler

import (
	"encoding/json"
	"fmt"
	"github.com/everpan/mdmg/pkg/base/entity"
	"github.com/everpan/mdmg/pkg/base/log"
	"github.com/everpan/mdmg/pkg/ctx"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"strings"
)

var logger = log.GetLogger()
var EntityGroupHandler = &ctx.IcGroupPathHandler{
	GroupPath: "/entity",
	Handlers: []*ctx.IcPathHandler{
		{
			Path:    "/meta/list/:page?", // page like 2-20, pageNum = 2, pageSize: 20
			Method:  fiber.MethodGet,
			Handler: metaList,
		},
		{
			Path:    "/meta/detail/:classNameOrId?",
			Method:  fiber.MethodGet,
			Handler: metaDetail,
		},
		{
			Path:    "/meta",
			Method:  fiber.MethodPost, // 新增
			Handler: metaCreate,
		},
		{
			Path:    "/meta",
			Method:  fiber.MethodPut, // modify
			Handler: metaDelete,
		},
		{
			Path:    "/meta",
			Method:  fiber.MethodPut, // modify
			Handler: metaDelete,
		},
	},
}

// metaDetail 通过class name or id获取元信息
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

// metaList 按照 ic_entity_class 列出与之相关的 ic_cluster_table
func metaList(c *ctx.IcContext) error {
	fc := c.FiberCtx()
	pageInfo := fc.Params("page")
	if pageInfo == "" {
		c.Page.Reset()
	} else {
		sp := strings.Split(pageInfo, "-")
		c.Page.PageNo, _ = strconv.Atoi(sp[0])
		c.Page.PageSize = 20
		if len(sp) > 1 {
			c.Page.PageSize, _ = strconv.Atoi(sp[1])
		}
	}
	// fmt.Printf("-- page %v\n", c.Page)
	count, _ := c.Engine().Count(&entity.IcEntityClass{TenantId: c.Tenant().Idx})
	offset := c.Page.CalCountOffset(int(count))

	var eClasses []*entity.IcEntityClass
	e := c.Engine().Limit(c.Page.PageSize, offset).Where("tenant_id = ?", c.Tenant().Idx).Find(&eClasses)
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

func metaCreate(c *ctx.IcContext) error {
	var (
		tenantId = c.Tenant().Idx
		fc       = c.FiberCtx()
	)
	meta, err := parseMetaFromBody(fc)
	if err != nil {
		return ctx.SendError(fc, fiber.StatusBadRequest, err)
	}
	/* 先增加class｜可以只有class，再增加cluster*/
	// [1] 校验，所有id必需为0; 同时将tenant id填入
	if meta.EntityClass.ClassId != 0 {
		err = fmt.Errorf("entity class id existed")
		return ctx.SendError(fc, fiber.StatusBadRequest, err)
	}
	if len(strings.TrimSpace(meta.EntityClass.ClassName)) == 0 {
		err = fmt.Errorf("entity class name is required")
		return ctx.SendError(fc, fiber.StatusBadRequest, err)
	}
	meta.EntityClass.TenantId = tenantId
	for _, table := range meta.ClusterTables {
		if table.ClusterId != 0 {
			err = fmt.Errorf("cluster table existed")
			return ctx.SendError(fc, fiber.StatusBadRequest, err)
		}
		if len(strings.TrimSpace(table.ClusterTableName)) == 0 {
			err = fmt.Errorf("cluster table name is required")
			return ctx.SendError(fc, fiber.StatusBadRequest, err)
		}
		table.TenantId = tenantId
	}
	// [2] register class
	_, err = c.EntityCtx().RegisterEntityClass(meta.EntityClass)
	if nil != err {
		return ctx.SendError(fc, fiber.StatusInternalServerError, err)
	}
	// [3] register cluster table
	for _, table := range meta.ClusterTables {
		table.ClassId = meta.EntityClass.ClassId
		err = c.EntityCtx().AddClusterTable(table)
		if err != nil {
			return ctx.SendError(fc, fiber.StatusInternalServerError, err)
		}
	}
	return ctx.SendSuccess(fc, meta)
}

func parseMetaFromBody(fc *fiber.Ctx) (meta *entity.IcEntityMeta, err error) {
	if len(fc.Body()) == 0 {
		return nil, fmt.Errorf("no body")
	}
	meta = &entity.IcEntityMeta{}
	err = json.Unmarshal(fc.Body(), meta)
	if nil != err || meta.EntityClass == nil {
		err = fmt.Errorf("body: %v, error: %v", string(fc.Body()), err)
		return nil, err
	}
	return meta, nil
}

func metaDelete(c *ctx.IcContext) error {
	fc := c.FiberCtx()
	meta, err := parseMetaFromBody(fc)
	if err != nil {
		return ctx.SendError(fc, fiber.StatusBadRequest, err)
	}
	if meta.EntityClass.ClassId != 0 {
		// del by class id
	}

	return nil
}
