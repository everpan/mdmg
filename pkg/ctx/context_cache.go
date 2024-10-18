package ctx

import (
	"github.com/everpan/mdmg/pkg/base/log"
	"github.com/everpan/mdmg/pkg/base/store"
	"github.com/everpan/mdmg/pkg/base/tenant"
	"github.com/everpan/mdmg/pkg/config"
	"github.com/gofiber/fiber/v2"
	"k8s.io/apimachinery/pkg/api/errors"
)

var (
	logger = log.GetLogger()
	cache  store.OneLevelMap[string, *IcContext]
)

func AcquireIcContextFromTenantId(fc *fiber.Ctx) (*IcContext, error) {
	tenantSid := fc.GetRespHeader("X-Tenant-Sid", tenant.TestTenantInfo.SId)
	//if tenantSid == tenant.TestSid {
	//	logger.Info("tenant id is the test id", zap.String("id", tenantSid))
	//}
	ctx, ok := cache.Get(tenantSid)
	if !ok {
		info, err := tenant.AcquireInfoBySid(tenantSid)
		if err != nil {
			e := errors.NewBadRequest("bad tenantSid:" + tenantSid + ",error:" + err.Error())
			return nil, e
		}
		if info == nil {
			e := errors.NewBadRequest("tenantSid:" + tenantSid + " not found")
			return nil, e
		}
		engine, err := config.AcquireEngine(info.Driver, info.ConnectString)
		if nil != err {
			e := errors.NewBadRequest("can not acquire engine for tenantSid:" + tenantSid + ",error:" + err.Error())
			return nil, e
		}
		ctx = NewContextWithParams(fc, info, nil, engine, "")
		cache.Set(tenantSid, ctx)
	}
	return ctx, nil
}
