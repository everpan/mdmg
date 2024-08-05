package icode

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	v8 "rogchap.com/v8go"
	"xorm.io/xorm"
)

type Ctx struct {
	fbCtx *fiber.Ctx
	V8Ctx *v8.Context
	Orm   *xorm.Engine
}

// 每个`fiber.Ctx`对应一个
var pool = map[*fiber.Ctx]*Ctx{}

func CreateV8Ctx(fb *fiber.Ctx) *v8.Context {
	iso := v8.NewIsolate()
	obj := v8.NewObjectTemplate(iso)
	obj.Set("icode", ExportObject(fb, iso))

	v8ctx := v8.NewContext(iso, obj)
	logger.Info("create v8 context", zap.Any("fbCtx", fb))
	return v8ctx
}

func createCtx(fb *fiber.Ctx) *Ctx {
	return &Ctx{
		fbCtx: fb,
		V8Ctx: CreateV8Ctx(fb),
	}
}

func GetCtx(fb *fiber.Ctx) *Ctx {
	ctx, ok := pool[fb]
	if !ok {
		ctx = createCtx(fb)
		pool[fb] = ctx
	}
	//logger.Info("get ctx",
	//	zap.Uintptr("fb", uintptr(unsafe.Pointer(fb))),
	//	zap.Uintptr("ctx", uintptr(unsafe.Pointer(ctx))))
	return ctx
}

func (fc *Ctx) Dispose() {
	fc.V8Ctx.Isolate().Dispose()
	fc.V8Ctx.Close()
}

func (fc *Ctx) SetFiberCtx(ctx *fiber.Ctx) {
	fc.fbCtx = ctx
}

func (fc *Ctx) GetFiberCtx() *fiber.Ctx {
	return fc.fbCtx
}

func DisposeCtxPool() {
	for _, ctx := range pool {
		ctx.Dispose()
		ctx.fbCtx = nil
		ctx.Orm = nil
	}
}
