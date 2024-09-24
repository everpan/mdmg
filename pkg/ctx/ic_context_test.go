package ctx

import (
	v8 "rogchap.com/v8go"
	"xorm.io/xorm"
)

func (c *IcContext) SetEngine(engine *xorm.Engine) {
	c.engine = engine
}
func (c *IcContext) SetV8Ctx(context *v8.Context) {
	c.v8Ctx = context
}
