package handler

import (
	"github.com/everpan/mdmg/pkg/base/tenant"
	"github.com/everpan/mdmg/pkg/ctx"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http/httptest"
	"strconv"
	"testing"
)

func Test_detail(t *testing.T) {
	tests := []struct {
		name    string
		id      int32
		wantErr string
	}{
		{"invalid id = 0", 0, "gt zero"},
	}
	app := fiber.New()
	ctx.AppRouterAddGroup(app, EntityGroupHandler)
	engine := CreateSeedDataSqlite3Engine("seed_data_test.db", false)
	tenant.SetSysEngine(engine)

	target := "/entity/meta/"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target2 := target + strconv.FormatInt(int64(tt.id), 10)
			req := httptest.NewRequest(fiber.MethodGet, target2, nil)
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
