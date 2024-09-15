package v8runtime

import (
	v8 "rogchap.com/v8go"
	"xorm.io/xorm"
)

func (c *Context) SetEngine(engine *xorm.Engine) {
	c.db = engine
}
func (c *Context) SetV8Ctx(context *v8.Context) {
	c.v8Ctx = context
}
