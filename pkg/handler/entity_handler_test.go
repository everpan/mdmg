package handler

import (
	"bytes"
	"encoding/json"
	"github.com/everpan/mdmg/pkg/base/entity"
	"github.com/everpan/mdmg/pkg/base/tenant"
	"github.com/everpan/mdmg/pkg/ctx"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"io"
	"math"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"xorm.io/xorm"
)

func Test_meta_detail(t *testing.T) {
	// ** 必需采用seed data 先构建相关表与数据
	tests := []struct {
		name      string
		id        int32
		className string
		wantErr   string
	}{
		{"invalid id empty", math.MaxInt32, "", "class name or id not specified"},
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

	target := "/entity/meta/detail/"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target2 := target + strconv.FormatInt(int64(tt.id), 10)
			if tt.id == math.MaxInt32 { // empty
				target2 = target
			}
			if tt.className != "" {
				target2 = target + tt.className
			}
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

func Test_meta_list(t *testing.T) {
	tests := []struct {
		name    string
		param   string
		retSize int
		want    string
	}{
		{"page 0", "", 20, "data\":[{\"entity_class\":{\"class_id\":1"},
		{"page 1", "/1-20", 20, "data\":[{\"entity_class\":{\"class_id\":1"},
		{"page 1", "/2-20", 20, "{\"code\":0,\"data\":[{\"entity_class\":{\"class_id\":21,"},
		{"page 3, left 11", "/3-20", 11, "page\":{\"page_size\":20,\"page_no\":3,\"page_count\":3,"},
		{"page 99, no data", "/99-20", 0, "data\":[]"},
		{"page 5, size 10 , left 1", "/6-10", 1, "{\"page_size\":10,\"page_no\":6,\"page_count\":6,\"record_count\":51}"},
	}
	app := fiber.New()
	ctx.AppRouterAddGroup(app, EntityGroupHandler)
	engine := CreateSeedDataSqlite3Engine("seed_data_test.db", false)
	tenant.SetSysEngine(engine)

	target := "/entity/meta/list"
	var target2 string
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target2 = target + tt.param
			req := httptest.NewRequest(fiber.MethodGet, target2, nil)
			resp, err := app.Test(req)
			if nil != err {
				assert.Contains(t, err.Error(), tt.want)
			}
			body, _ := io.ReadAll(resp.Body)
			// t.Log(string(body))
			assert.Contains(t, string(body), tt.want)

			var r = ctx.ICodeResponse{}
			e := json.Unmarshal(body, &r)
			assert.Nil(t, e)
			//随着数据的增加，可能判定条件发生变化，注意观察；此判断不稳定
			assert.GreaterOrEqual(t, tt.retSize, len(r.Data.([]any)))
		})
	}
}

func clearData(engine *xorm.Engine) {
	sql := "delete from ic_entity_class where class_name in ('only_entity_class_test' ,'entity_class_test','less_fields')"
	engine.Exec(sql)
	sql = "delete from ic_cluster_table where cluster_table_name in ('cluster_table_name','less_fields_cluster')"
	engine.Exec(sql)
}

func Test_metaAdd(t *testing.T) {
	tests := []struct {
		name string
		meta *entity.IcEntityMeta
		body string //优先
		want string
	}{
		{"body is null", nil, "null", "body: null"},
		{"body fmt error", nil, "bad fmt", "invalid character"},
		{"no body", nil, " ", "{\"code\":-1,\"message\":\"no body\"}"},
		{"only entity class", &entity.IcEntityMeta{EntityClass: &entity.IcEntityClass{ClassName: "only_entity_class_test", PkColumn: "idx"}},
			"", "\"tenant_id\":1"},
		{"same entity class throw constraint", &entity.IcEntityMeta{EntityClass: &entity.IcEntityClass{ClassName: "only_entity_class_test", PkColumn: "idx"}},
			"", "UNIQUE constraint"},
		{"with cluster table", &entity.IcEntityMeta{EntityClass: &entity.IcEntityClass{ClassName: "entity_class_test", PkColumn: "idx"},
			ClusterTables: []*entity.IcClusterTable{{ClusterTableName: "cluster_table_name"}}},
			"", "\"code\":0"},
		{"less fields in json", nil,
			`{"entity_class":{ "name":"less_fields","class_desc":"","pk_column":"idx"},"cluster_tables":[{"table":"less_fields_cluster"}]}`,
			`{"code":0,"data":`},
	}

	app := fiber.New()
	ctx.AppRouterAddGroup(app, EntityGroupHandler)
	engine := CreateSeedDataSqlite3Engine("seed_data_test.db", false)
	clearData(engine)
	tenant.SetSysEngine(engine)
	target := "/entity/meta/"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			var buf bytes.Buffer
			if len(tt.body) > 0 {
				buf.WriteString(strings.TrimSpace(tt.body))
			} else {
				json.NewEncoder(&buf).Encode(tt.meta)
			}
			req = httptest.NewRequest(fiber.MethodPost, target, &buf)
			resp, err := app.Test(req)
			if err != nil {
				assert.Contains(t, err.Error(), tt.want)
			}
			body, _ := io.ReadAll(resp.Body)
			// t.Log(string(body))
			assert.Contains(t, string(body), tt.want)
			t.Log(string(body))
			var r = ctx.ICodeResponse{}
			e := json.Unmarshal(body, &r)
			assert.Nil(t, e)
			// assert.Equal(t, tt.retSize, len(r.Data.([]any)))
		})
	}
}
