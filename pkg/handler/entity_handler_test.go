package handler

import (
	"github.com/everpan/mdmg/pkg/base/tenant"
	"github.com/everpan/mdmg/pkg/ctx"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"io"
	"math"
	"net/http/httptest"
	"strconv"
	"testing"
)

func Test_detail(t *testing.T) {
	// ** 必需采用seed data 先构建相关表与数据
	tests := []struct {
		name      string
		id        int32
		className string
		wantErr   string
	}{
		{"invalid id empty", math.MaxInt32, "", "class not specified"},
		{"invalid id 0", 0, "", "gt zero"},
		{"invalid id 99", 99, "", "not found"},
		{"invalid id 1", 1, "", "{\"code\":0,\"data\":{\"entity_class\":{\"class_id\":1"},
		{"class name: user_not_exist", 0, "user_not_exist", "className:user_not_exist tenantId:1 not found"},
		{"class name: user", 0, "user", "\"cluster_tables\":[{\"class_id"},
	}
	app := fiber.New()
	ctx.AppRouterAddGroup(app, EntityGroupHandler)
	engine := CreateSeedDataSqlite3Engine("seed_data_test.db", false)
	tenant.SetSysEngine(engine)

	target := "/entity/meta/"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target2 := target + strconv.FormatInt(int64(tt.id), 10)
			if tt.id == math.MaxInt32 { // empty
				target2 = target
			}
			if tt.className != "" {
				target2 = target + tt.className
			}
			req := httptest.NewRequest(fiber.MethodPost, target2, nil)
			resp, err := app.Test(req)
			if nil != err {
				assert.Contains(t, err.Error(), tt.wantErr)
			}
			body, _ := io.ReadAll(resp.Body)
			// t.Log(string(body))
			assert.Contains(t, string(body), tt.wantErr)
		})
	}
}
