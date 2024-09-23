package entity

import (
	"fmt"
	"github.com/everpan/mdmg/dsl"
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

var (
	primaryTable = &dsl.Meta{Table: dsl.Table{Name: "a"}, Columns: []*dsl.Column{{Name: "a1"}, {Name: "idx", IsPrimaryKey: true}}}
	data         = []*dsl.Meta{
		{Table: dsl.Table{Name: "a"}, Columns: []*dsl.Column{{Name: "a1"}}},
		{Table: dsl.Table{Name: "b"}, Columns: []*dsl.Column{{Name: "a1"}, {Name: "a2"}}},
		{Table: dsl.Table{Name: "c"}, Columns: []*dsl.Column{{Name: "a1"}, {Name: "a2"}, {Name: "a3"}}},
	}
)

func TestGenerateLeftJoinSQL(t *testing.T) {
	type args struct {
		tables []*dsl.Meta
		key    string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{"zero table left join", args{nil, "key"},
			""},
		{"one table left join", args{data[0:1], "key"},
			"\nleft join a as t1 on t0.key=t1.key"},
		{"two tables left join", args{data[0:2], "key"},
			"\nleft join a as t1 on t0.key=t1.key\nleft join b as t2 on t0.key=t2.key"},
		{"three tables left join", args{data[0:3], "key"},
			"\nleft join a as t1 on t0.key=t1.key\nleft join b as t2 on t0.key=t2.key\nleft join c as t3 on t0.key=t3.key"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, GenerateLeftJoinSQL(tt.args.tables, tt.args.key), "GenerateLeftJoinSQL(%v, %v)", tt.args.tables, tt.args.key)
		})
	}
}

func Test_filterRepeatedColumns(t *testing.T) {
	type args struct {
		clusterTables []*dsl.Meta
	}
	tests := []struct {
		name  string
		args  args
		count int
	}{
		{"empty", args{nil}, 0},
		{"one column", args{data[0:2]}, 1},
		{"one column", args{data}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := filterRepeatedColumns(nil, tt.args.clusterTables)
			assert.Equalf(t, tt.count, len(r), "filterRepeatedColumns(%v)", tt.args.clusterTables)
		})
	}
}

func TestGenerateSelectColumnsSQL(t *testing.T) {
	tests := []struct {
		name          string
		primaryTables *dsl.Meta
		clusterTables []*dsl.Meta
		want          string
		wantErr       string
	}{
		{"no primary key", data[0], data, "", "not fount primary key in table"},
		{"2nd table aliased", primaryTable, data[0:1],
			"select t0.idx, t0.a1\n, t1.a1 as a1_1", ""},
		{"3rd table aliased", primaryTable, data[0:2],
			"select t0.idx, t0.a1\n, t1.a1 as a1_1\n, t2.a1 as a1_2, t2.a2", ""},
		{"4th table aliased", primaryTable, data[0:3],
			"select t0.idx, t0.a1\n, t1.a1 as a1_1\n, t2.a1 as a1_2, t2.a2 as a2_2\n, t3.a1 as a1_3, t3.a2 as a2_3, t3.a3", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateSelectColumnsSQL(tt.primaryTables, tt.clusterTables)
			assert.Equalf(t, tt.want, got, "GenerateSelectColumnSQL(%v, %v)", tt.primaryTables, tt.clusterTables)
			if nil != err && tt.wantErr != "" {
				assert.Containsf(t, err.Error(), tt.wantErr, "GenerateSelectColumnSQL(%v, %v)", tt.primaryTables, tt.clusterTables)
			}
		})
	}
}
