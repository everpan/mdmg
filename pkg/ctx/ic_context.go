package ctx

import (
	"fmt"
	"github.com/everpan/mdmg/pkg/base/entity"
	"github.com/everpan/mdmg/pkg/base/tenant"
	"github.com/everpan/mdmg/pkg/base/v8runtime"
	"github.com/everpan/mdmg/utils"
	"github.com/gofiber/fiber/v2"
	v8 "rogchap.com/v8go"
	"xorm.io/xorm"
)

type IcPathHandler struct {
	Path    string
	Method  string
	Handler IcHandler
}

type IcGroupPathHandler struct {
	GroupPath string
	Handlers  []*IcPathHandler
}

type IcHandler func(ctx *IcContext) error

func (h IcHandler) WrapHandler() fiber.Handler {
	return func(fc *fiber.Ctx) error {
		ctx, err := AcquireIcContextFromTenantId(fc)
		if err != nil || ctx == nil {
			if ctx == nil {
				err = fmt.Errorf("ctx is nil,%v", err.Error())
			}
			return SendError(fc, fiber.StatusBadRequest, err)
		}
		ctx.fc = fc
		return h(ctx)
	}
}

type IcPage struct {
	Size   int `json:"size"`
	Number int `json:"pageNumber"`
	// Where      string // where id > 10 效率
}

func NewIcPage() *IcPage {
	p := &IcPage{}
	p.Reset()
	return p
}

func (p *IcPage) Reset() {
	p.Size = 20
	p.Number = 0
}

type IcContext struct {
	fc            *fiber.Ctx
	tenant        *tenant.IcTenantInfo
	v8Ctx         *v8.Context
	engine        *xorm.Engine
	moduleVersion string
	EntityCtx     *entity.Context
	Page          *IcPage
}

func NewContextWithParams(
	fc *fiber.Ctx,
	tenant *tenant.IcTenantInfo,
	v8Ctx *v8.Context,
	engine *xorm.Engine,
	moduleVersion string,
) *IcContext {
	ctx := &IcContext{
		fc:            fc,
		tenant:        tenant,
		v8Ctx:         v8Ctx,
		engine:        engine,
		moduleVersion: moduleVersion,
	}
	if ctx.v8Ctx == nil {
		ctx.CreateV8Context()
	}
	if ctx.EntityCtx == nil {
		ctx.EntityCtx = entity.NewContext(engine, tenant.Idx)
	}
	if ctx.Page == nil {
		ctx.Page = NewIcPage()
	}
	return ctx
}
func (c *IcContext) FiberCtx() *fiber.Ctx {
	return c.fc
}
func (c *IcContext) V8Ctx() *v8.Context {
	return c.v8Ctx
}
func (c *IcContext) Engine() *xorm.Engine {
	return c.engine
}
func (c *IcContext) ModuleVersion() string {
	return c.moduleVersion
}
func (c *IcContext) SetModuleVersion(moduleVersion string) {
	c.moduleVersion = moduleVersion
}
func (c *IcContext) RunScript(source string, origin string) (*v8.Value, error) {
	return c.v8Ctx.RunScript(source, origin)
}

func (c *IcContext) RunScriptRetAny(source string, origin string) (any, *v8.Value, error) {
	v, err := c.v8Ctx.RunScript(source, origin)
	if err != nil {
		return nil, nil, err
	}
	g, e := utils.ToGoValue(c.v8Ctx, v)
	return g, v, e
}

func (c *IcContext) CreateV8Context() *v8.Context {
	iso := v8.NewIsolate()
	icObj := v8.NewObjectTemplate(iso)
	obj := v8.NewObjectTemplate(iso)
	ctxObj := c.ExportV8ObjectTemplate(iso)
	_ = obj.Set("ctx", ctxObj)
	_ = obj.Set("engine", v8runtime.ExportXormObject(c.engine, iso))
	_ = icObj.Set("__ic", obj)
	v8ctx := v8.NewContext(iso, icObj)
	// icode.logger.IcTenantInfo("create v8 context", zap.Any("fbCtx", fb))
	c.v8Ctx = v8ctx
	return v8ctx
}

func (c *IcContext) ExportV8ObjectTemplate(iso *v8.Isolate) *v8.ObjectTemplate {
	mf := v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		m, _ := utils.SplitModuleVersion(c.moduleVersion)
		jv, _ := v8.NewValue(iso, m)
		return jv
	})
	vf := v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		_, v := utils.SplitModuleVersion(c.moduleVersion)
		jv, _ := v8.NewValue(iso, v)
		return jv
	})
	ti := v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		jv, _ := utils.ToJsValue(info.Context(), c.tenant)
		return jv
	})
	ctxObj := v8runtime.ExportObject(c.fc, iso)
	_ = ctxObj.Set("module", mf)
	_ = ctxObj.Set("version", vf)
	_ = ctxObj.Set("tenant", ti)
	return ctxObj
}
