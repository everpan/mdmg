package handler

import (
	"github.com/everpan/mdmg/pkg/store"
	"github.com/everpan/mdmg/pkg/tenant"
	"github.com/gofiber/fiber/v2"
	"k8s.io/apimachinery/pkg/api/errors"
)

var (
	cache store.OneLevelMap[string, *Context]
)

func AcquireContext(fc *fiber.Ctx) (*Context, error) {
	tenantSid := fc.GetRespHeader("X-Tenant-Sid", tenant.DefaultGuidNamespace)
	ctx, ok := cache.Get(tenantSid)
	if !ok {
		info, err := tenant.AcquireTenantInfoBySid(tenantSid)
		if err != nil {
			e := errors.NewBadRequest("bad tenantSid:" + tenantSid + ",error:" + err.Error())
			return nil, e
		}
		if info == nil {
			e := errors.NewBadRequest("tenantSid:" + tenantSid + " not found")
			return nil, e
		}
		engine, err := tenant.AcquireEngine(info)
		if nil != err {
			e := errors.NewBadRequest("can not acquire engine for tenantSid:" + tenantSid + ",error:" + err.Error())
			return nil, e
		}
		ctx = &Context{
			fc:     fc,
			tenant: info,
			db:     engine,
		}
		ctx.CreateV8Context()
		cache.Set(tenantSid, ctx)
	}
	return ctx, nil
}
