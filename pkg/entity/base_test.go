package entity

import (
	"fmt"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"strings"
	"xorm.io/xorm"
)

func setUp() {
	engine, err := xorm.NewEngine("sqlite3", "./entity_test.db")
	if err != nil {
		fmt.Printf("failed to setup test engine: %v", err.Error())
	}
	engine.ShowSQL(true)
	SetEngine(engine)
	InitTable()
}

func tearDown() {
	engine.Close()
	// os.Remove("./entity_test.db")
}

func TestGetEntityClass(t *testing.T) {
	setUp()
	emptyErrFn := func(testingT assert.TestingT, err error, i ...interface{}) bool {
		return false
	}
	entityClassCache.Set(100, &IcEntityClass{ClassId: 100})
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
			got, err := GetEntityClass(tt.classId)
			if !tt.wantErr(t, err, fmt.Sprintf("GetEntityClass(%v)", tt.classId)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetEntityClass(%v)", tt.classId)
		})
	}
}
