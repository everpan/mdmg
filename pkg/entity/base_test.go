package entity

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"strings"
	"xorm.io/xorm"
)

func setUp(tName string) *Context {
	dbName := fmt.Sprintf("./%s_test.db", tName)
	_ = os.Remove(dbName)

	engine, err := xorm.NewEngine("sqlite3", dbName)
	if err != nil {
		fmt.Printf("failed to setup test engine: %v", err.Error())
	}
	engine.ShowSQL(true)
	ctx := NewContext(engine)
	ctx.InitTable()
	return ctx
}

var emptyErrFn = func(testingT assert.TestingT, err error, i ...interface{}) bool {
	return false
}

func TestGetEntityClass(t *testing.T) {
	ctx := setUp("entity_class")

	ctx.entityClassCache.Set(100, &IcEntityClass{ClassId: 100})
	tests := []struct {
		name    string
		classId uint32
		want    *IcEntityClass
		wantErr assert.ErrorAssertionFunc
	}{
		{"fetch nil", uint32(2000), nil,
			func(testingT assert.TestingT, err error, i ...interface{}) bool {
				b := strings.Contains(err.Error(), "2000 not found")
				if !b {
					t.Error("not contains `2000 not found`")
				}
				return b
			},
		},
		{"fetch 0", uint32(0), nil,
			func(testingT assert.TestingT, err error, i ...interface{}) bool {
				//t.Error(err)
				if err.Error() != "classId is 0" {
					t.Error("want err.Error == classId is 0, but err is ", err)
					return false
				}
				return true
			},
		},
		{"fetch exist 100 ", uint32(100), &IcEntityClass{ClassId: 100}, emptyErrFn},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ctx.GetEntityClass(tt.classId)
			if !tt.wantErr(t, err, fmt.Sprintf("GetEntityClass(%v)", tt.classId)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetEntityClass(%v)", tt.classId)
		})
	}
}

func TestRegisterEntityClass(t *testing.T) {
	// tearDown()
	ctx := setUp("register_entity_class")
	ec := &IcEntityClass{0, "user", "user info 1", "user_id"}
	e2 := *ec
	e2Want := e2
	e2Want.ClassId = 1
	// clear
	_, err := ctx.engine.Exec("drop table ic_entity_class")
	ctx.InitTable()
	if err != nil {
		t.Error(err)
	}
	tests := []struct {
		name    string
		ec      *IcEntityClass
		want    *IcEntityClass
		wantErr assert.ErrorAssertionFunc
	}{
		{"insert 0", ec, &e2Want, emptyErrFn},
		{"insert same unique class name", &e2, nil,
			func(t assert.TestingT, err error, i ...interface{}) bool {
				b := strings.Contains(err.Error(), "UNIQUE constraint failed")
				if !b {
					t.Errorf(err.Error())
				}
				return b
			},
		},
	}
	for _, tt := range tests {
		//t.Run(tt.name, func(t *testing.T) {
		t.Log(tt.ec)
		got, err := ctx.RegisterEntityClass(tt.ec)
		if err != nil && !tt.wantErr(t, err, fmt.Sprintf("RegisterEntityClass(%v)", tt.ec)) {
			t.Error(err.Error())
		}
		assert.Equalf(t, tt.want, got, "RegisterEntityClass(%v)", tt.ec)
		//})
	}
	e, err := ctx.GetEntityClassByName(ec.ClassName)
	if err != nil {
		t.Error(err)
	}
	assert.Equalf(t, e, ec, "GetEntityClassByName(%s)", ec.ClassName)
}
