package v8runtime

import (
	"net/http/httptest"
	"testing"

	"github.com/everpan/mdmg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	v8 "rogchap.com/v8go"
)

func TestExportObject(t *testing.T) {
	tests := []struct {
		name   string
		path   string
		target string
		script string
		want   func(ctx *IcContext, value *v8.Value) bool
	}{
		{"undefined", "", "/", "", func(ctx *IcContext, value *v8.Value) bool {
			// logger.Info("run", zap.Any("val", value))
			return value.String() == "undefined"
		}},
		{"global object", "/test/:p1/:p2/*", "/test/a/:b/c/d",
			`
let icode = __ic.ctx
(() => {
let accept = icode.header().Accept
return {
	code: 0,
	data: {
		sql: "select * from user"
	},
	header: icode.header("content-type"),
	headers: icode.header(),
	query: icode.query('key'),
	queries: icode.query(),
	param: icode.param("module"),
	params: icode.param(),
	accept,
	base: icode.baseURL(),
	originURL: icode.originURL()
}
})()`, func(ctx *IcContext, value *v8.Value) bool {
				gv, _ := utils.ToGoValue(ctx.V8Ctx(), value)
				// logger.Info("run", zap.Any("val", gv), zap.String("type", reflect.TypeOf(gv).String()))
				jv0 := gv.(map[string]interface{})
				params := jv0["params"].(map[string]any)
				// logger.Info("params", zap.Any("params", params), zap.String("type", reflect.TypeOf(jv0["params"]).String()))
				return params["p1"] == "a" &&
					params["*1"] == "c/d" &&
					jv0["data"].(map[string]any)["sql"] == "select * from user"
			}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			app.Get(tt.path, func(c *fiber.Ctx) error {
				ctx := NewContextWithParams(c, nil, nil, nil, "")
				val, err := ctx.RunScript(tt.script, "test.js")
				assert.Nil(t, err)
				assert.True(t, tt.want(ctx, val))
				return nil
			})
			resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, tt.target, nil), -1)
			assert.Nil(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, resp.StatusCode, 200)
		})
	}
}
