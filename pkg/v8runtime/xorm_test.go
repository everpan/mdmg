package v8runtime

import (
	"fmt"
	"github.com/everpan/mdmg/utils"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"os"
	v8 "rogchap.com/v8go"
	"testing"
	"time"
	xr "xorm.io/xorm"
	"xorm.io/xorm/log"
)

var _ctx = &Context{}
var testDbFilename = "xorm_test.db"

func setup() {
	os.Remove(testDbFilename)
	var err error
	eng, err := xr.NewEngine("sqlite3", testDbFilename)
	eng.ShowSQL(true)
	eng.Logger().SetLevel(log.LOG_DEBUG)
	if err != nil {
		panic(err)
	}
	// init v8ctx
	// _ctx := &Ctx{db: eng}
	_ctx.SetEngine(eng)
	iso := v8.NewIsolate()
	obj := ExportXormObject(_ctx.db, iso)
	_ctx.SetV8Ctx(v8.NewContext(iso, obj))

	buildTestData()
}

func teardown() {
	_ctx.Engine().Close()
	_ctx.V8Ctx().Close()
}

func buildTestData() {
	t := time.Now()
	type User struct {
		ID        int64     `xorm:"id pk autoincr"`
		Name      string    `xorm:"name UNIQUE NOT NULL"`
		Age       int       `xorm:"age UNIQUE NOT NULL"`
		Birthday  time.Time `xorm:"birthday TIMESTAMP"`
		CreatedAt time.Time `xorm:"create_at TIMESTAMP"`
	}
	var user = []User{
		{0, "name1", 23, t.Add(1), t},
		{0, "name2", 24, t.Add(2), t},
		{0, "name3", 24, t.Add(2), t},
	}
	err := _ctx.Engine().CreateTables(user[0])
	if err != nil {
		panic(err)
	}
	_, err = _ctx.Engine().Insert(user)
	if err != nil {
		panic(err)
	}
	// fmt.Printf("after insert:%v\n", user)
}

func TestQuery(t *testing.T) {
	tests := []struct {
		name string
		sql  string
		want func(t *testing.T, v *v8.Value, g any)
	}{
		{"select all", "select * from user", func(t *testing.T, v *v8.Value, g any) {
			//d, _ := json.MarshalIndent(g, "", " ")
			//t.Log(string(d))
			assert.GreaterOrEqual(t, 3, len(g.([]any)))
		}},
		{"select one", "select * from user where id=1", func(t *testing.T, v *v8.Value, g any) {
			assert.Equal(t, 1, len(g.([]any)))
		}},
		{"select zero", "select * from user where id=0", func(t *testing.T, v *v8.Value, g any) {
			assert.Equal(t, "null", v.String())
			assert.Nil(t, g)
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := _ctx.RunScript(fmt.Sprintf("query(\"%s\")", tt.sql), "xorm_test.js")
			if err != nil {
				t.Fatal(err)
			}
			g, _ := utils.ToGoValue(_ctx.V8Ctx(), r)
			tt.want(t, r, g)
		})
	}
}

func TestExec(t *testing.T) {
	tests := []struct {
		name string
		exec string
		want func(t *testing.T, v *v8.Value, g any)
	}{
		{"insert one",
			`exec("insert into user(name,age) values (?,89)","my name")`,
			func(t *testing.T, v *v8.Value, g any) {
				v2, _, err := _ctx.RunScriptRetAny("query('select * from user where name=\"my name\"')", "xorm_test.js")
				assert.NoError(t, err)
				assert.Equal(t, 1, len(v2.([]any)))
				v3, _, err := _ctx.RunScriptRetAny("query('select count(1) as c from user')", "xorm_test.js")
				assert.NoError(t, err)
				// todo why return type is string
				// t.Log("type", reflect.TypeOf(v3.([]any)[0].(map[string]any)["c"]).String())
				assert.Equal(t, "4", v3.([]any)[0].(map[string]any)["c"])
			},
		},
		{"delete", `exec("delete from user where id=?",1)`, func(t *testing.T, v *v8.Value, g any) {
			v2, _, err := _ctx.RunScriptRetAny("query('select * from user where id=?',1)", "xorm_test.js")
			assert.NoError(t, err)
			assert.Nil(t, v2)
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := _ctx.RunScript(tt.exec, "xorm_test.js")
			if err != nil {
				t.Fatal(err)
			}
			g, _ := utils.ToGoValue(_ctx.V8Ctx(), r)
			tt.want(t, r, g)
		})
	}
}

func TestMain(m *testing.M) {
	setup()
	if m.Run() == 0 {
		teardown()
		os.Remove(testDbFilename)
	} else {
		teardown()
	}
}
