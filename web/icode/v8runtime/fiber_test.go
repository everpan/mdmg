package v8runtime

import (
	"github.com/everpan/mdmg/utils"
	"github.com/everpan/mdmg/web/icode"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	v8 "rogchap.com/v8go"
	"testing"
)

func init() {
	icode.InitLogger()
}

func TestExportObject(t *testing.T) {
	tests := []struct {
		name   string
		path   string
		target string
		script string
		want   func(ctx *icode.Ctx, value *v8.Value) bool
	}{
		{"undefined", "", "/", "", func(ctx *icode.Ctx, value *v8.Value) bool {
			// logger.Info("run", zap.Any("val", value))
			return value.String() == "undefined"
		}},
		{"global object", "/test/:p1/:p2/*", "/test/a/:b/c/d",
			`(() => {
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
})()`, func(ctx *icode.Ctx, value *v8.Value) bool {
				gv, _ := utils.ToGoValue(ctx.V8Context(), value)
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
			var ctx *icode.Ctx
			app := fiber.New()
			defer app.Shutdown()
			defer icode.DisposeCtxPool()
			app.Get(tt.path, func(c *fiber.Ctx) error {
				ctx = icode.AcquireCtx(c)

				val, err := ctx.RunScript(tt.script, "test.js")
				assert.Nil(t, err)
				assert.True(t, tt.want(ctx, val))
				return nil
			})
			resp, _ := app.Test(httptest.NewRequest(fiber.MethodGet, tt.target, nil))
			assert.Equal(t, resp.StatusCode, 200)
		})
	}
}
