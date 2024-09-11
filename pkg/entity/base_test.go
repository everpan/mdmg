package entity

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"testing"
	"xorm.io/xorm"
)

func setUp() {
	engine, err := xorm.NewEngine("sqlite3", "./entity_test.db")
	if err != nil {
		fmt.Errorf("failed to setup test engine: %v", err)
	}
	engine.ShowSQL(true)
	SetEngine(engine)
	InitTable()
}

func tearDown() {
	engine.Close()
	// os.Remove("./entity_test.db")
}

func TestAcquireEntityBase(t *testing.T) {
	setUp()
	defer tearDown()
	tests := []struct {
		name string
		base []*IcEntityBase
	}{
		{"insert", []*IcEntityBase{
			{uint64(0), uint32(1), "key1"},
			{uint64(0), uint32(1), "key2"},
			{uint64(0), uint32(2), "key1"},
			{uint64(0), uint32(2), "key2"},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, b := range tt.base {
				AcquireEntityBase(b)
				assert.Greater(t, b.EntityId, uint64(0))
			}
		})
	}
}

func TestInsertUniqueEntityBase(t *testing.T) {
	setUp()
	defer tearDown()
	tests := []struct {
		name string
		base []*IcEntityBase
	}{
		{"insert new", []*IcEntityBase{
			{uint64(0), uint32(1), "key1"},
			{uint64(0), uint32(1), "key2"},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, b := range tt.base {
				AcquireEntityBase(b)
				assert.Greater(t, b.EntityId, uint64(0))
			}
		})
	}
}

func TestQueryEntityBaseByClsIdAndUKey(t *testing.T) {
	setUp()
	type args struct {
		clsId uint32
		ukey  string
	}
	tests := []struct {
		name string
		args args
		want *IcEntityBase
	}{
		{"fetch exist", args{uint32(1), "key1"},
			&IcEntityBase{EntityId: 1, ClassId: 1, EntityUKey: "key1"},
		},
		{"fetch not exist", args{uint32(9), "key1"}, nil},
		{"fetch not exist", args{uint32(1), "key9"}, nil},
		{"fetch not exist", args{uint32(0), "key9"}, nil},
		{"** fetch when cls=0 and ukey is exist", args{uint32(0), "key1"}, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, QueryEntityBaseByClsIdAndUKey(tt.args.clsId, tt.args.ukey),
				"QueryEntityBaseByClsIdAndUKey(%v, %v)", tt.args.clsId, tt.args.ukey)
		})
	}
}
