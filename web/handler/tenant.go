package handler

import (
	"github.com/everpan/mdmg/pkg/store"
	"github.com/everpan/mdmg/pkg/tenant"
	"github.com/gofiber/fiber/v2"
)

var (
	cache store.OneLevelMap[string, *Context]
)

func AcquireContext(fc *fiber.Ctx) *Context {
	tenantSid := fc.GetRespHeader("X-Tenant-Sid", tenant.DefaultGuidNamespace)
	ctx, ok := cache.Get(tenantSid)
	if !ok {
		info, _ := tenant.AcquireTenantInfoBySid(tenantSid)
		engine, _ := tenant.AcquireEngine(info)

		ctx = &Context{
			tenant: info,
			db:     engine,
		}
		ctx.CreateV8Context()
		cache.Set(tenantSid, ctx)
	}
	ctx.fc = fc
	return ctx
}
