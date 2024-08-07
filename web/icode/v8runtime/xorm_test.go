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
	"xorm.io/xorm"
	"xorm.io/xorm/log"
)

var eng *xorm.Engine
var v8ctx *v8.Context

func init() {
	fName := "xorm_test.db"
	os.Remove(fName)
	var err error
	eng, err = xorm.NewEngine("sqlite3", fName)
	eng.ShowSQL(true)
	eng.Logger().SetLevel(log.LOG_DEBUG)
	if err != nil {
		panic(err)
	}
	// init v8ctx
	ctx := &Ctx{db: eng}
	iso := v8.NewIsolate()
	obj := ExportXormObject(ctx, iso)
	v8ctx = v8.NewContext(iso, obj)

	buildTestData()
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
	err := eng.CreateTables(user[0])
	if err != nil {
		panic(err)
	}
	_, err = eng.Insert(user)
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
			r, err := v8ctx.RunScript(fmt.Sprintf("query(\"%s\")", tt.sql), "xorm_test.js")
			if err != nil {
				t.Fatal(err)
			}
			g, _ := utils.ToGoValue(v8ctx, r)
			tt.want(t, r, g)
		})
	}
}

func TestExecInsert(t *testing.T) {

}
