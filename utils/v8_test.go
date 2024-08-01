package utils

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
	"testing"
	"xorm.io/xorm"
)
import v8 "rogchap.com/v8go"

func TestOneVmMultiContext(t *testing.T) {
	iso := v8.NewIsolate()     // creates a new JavaScript VM
	ctx1 := v8.NewContext(iso) // new context within the VM
	ctx1.RunScript("const multiply = (a, b) => a * b", "math.js")
	v1, err := ctx1.RunScript("const result = multiply(2,3)", "math.js")
	if err != nil {
		t.Fatal(err)
	}
	v2, _ := ctx1.RunScript("result", "math.js") // return a value in JavaScript back to Go
	t.Logf("vaule: %v, %v", v1, v2)

	//
	ctx2 := v8.NewContext(iso) // another context on the same VM
	if v, err := ctx2.RunScript("multiply(3, 4)", "main.js"); err != nil {
		// this will error as multiply is not defined in this context
		require.Error(t, err)
	} else {
		t.Log(v)
	}
}

func TestUpdateJsObject(t *testing.T) {
	iso := v8.NewIsolate()
	defer iso.Dispose()
	eng, err := xorm.NewEngine("mysql", "devuser:devuser.COM2019@tcp(devmysql01.wiz.top:6033)/wiz_hr2?charset=utf8")
	if err != nil {
		t.Error(err)
	}

	tmpl := v8.NewObjectTemplate(iso)
	tmpl.Set("engine", eng)

	ctx := v8.NewContext(iso, tmpl)
	defer ctx.Close()

	v, err := ctx.RunScript("engine", "main.js")
	if err != nil {
		t.Error(err)
	}
	t.Log(v)
	///
	val, _ := v8.JSONParse(ctx, `{
		"a": 1,
		"b": "foo"
	}`)
	v, err = ctx.RunScript("jval.b", "main.js")
	t.Log(v)
	if err != nil {
		t.Error(err)
	}

	t.Logf("json val: %v %v", val, val.String())
	ctx.Global().Set("jval", val)
	v, err = ctx.RunScript("jval.b", "main.js")
	t.Log(v)
	if err != nil {
		t.Error(err)
	}
	// js function
	v, err = ctx.RunScript("fun  = p => {console.log(p)}; typeof(fun)", "main.js")
	t.Log(v)
	if err != nil {
		t.Error(err)
	}

}
