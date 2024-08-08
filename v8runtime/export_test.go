package v8runtime

import (
	"github.com/gofiber/fiber/v2"
	v8 "rogchap.com/v8go"
	"xorm.io/xorm"
)

func (c *Ctx) SetEngine(eng *xorm.Engine) {
	c.db = eng
}

func (c *Ctx) SetFiberCtx(ctx *fiber.Ctx) {
	c.fbCtx = ctx
}

func (c *Ctx) SetV8Ctx(v8Ctx *v8.Context) {
	c.v8Ctx = v8Ctx
}
