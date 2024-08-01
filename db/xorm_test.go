package db_test

import (
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"reflect"
	"testing"
	"xorm.io/xorm"
)

func TestDBMeta(t *testing.T) {
	var engine *xorm.Engine
	engine, err := xorm.NewEngine("mysql", "root:123456@tcp(127.0.0.1:3306)/wiz_hr2?charset=utf8")
	if err != nil {
		t.Error(err)
	}
	defer engine.Close()
	//engine.ShowSQL(true)
	//engine.Logger().SetLevel(log.LOG_DEBUG)

	if engine == nil {
		t.Fatalf("engine is nil")
	}
	metas, err := engine.DBMetas()
	t.Logf("metas.len:%d", len(metas))
	if err == nil {
		//for n, meta := range metas {
		//	j, e := json.Marshal(meta)
		//	if e != nil {
		//		t.Error(e)
		//	} else {
		//		t.Logf("%d %v\n", n, string(j))
		//	}
		//}
		//for n := 0; n < len(metas); n++ {
		//	j, e := json.Marshal(metas[n])
		//	if e != nil {
		//		t.Error(e)
		//	} else {
		//		t.Logf("%d %v\n", n, string(j))
		//	}
		//}
		//for n, meta := range metas {
		//	j, _ := json.Marshal(meta)
		//	t.Logf("%d %v\n", n, string(j))
		//}
	}
	// 获取columns信息
	j, _ := json.Marshal(metas[1])
	t.Logf("table:%s", string(j))
	table := metas[1]
	cols := table.Columns()
	for n, col := range cols {
		j, _ := json.Marshal(col)
		t.Logf("%d col:%s", n, string(j))
		break
	}
}

// TestDefer 测试defer的范围与时机
func TestDefer(t *testing.T) {
	r := make([]string, 0)
	func(fn func()) {
		r = append(r, "defer fun begin")
		defer fn() // must be at end
		r = append(r, "defer fun end")
	}(func() {
		r = append(r, "defer fun exec")
	})
	r = append(r, "end")
	want := []string{
		"defer fun begin", "defer fun end", "defer fun exec", "end",
	}
	if !reflect.DeepEqual(r, want) {
		t.Errorf("%v\n%v", r, want)
	}
}
