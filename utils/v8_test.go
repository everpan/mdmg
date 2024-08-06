package utils

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
	"testing"
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
