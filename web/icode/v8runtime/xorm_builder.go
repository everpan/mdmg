package v8runtime

import (
	v8 "rogchap.com/v8go"
	"xorm.io/builder"
)

func ExportXormBuilder(iso *v8.Isolate) *v8.ObjectTemplate {
	obj := v8.NewObjectTemplate(iso)
	return obj
}
func select_(iso *v8.Isolate) *v8.FunctionTemplate {
	return v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		info.This()
		builder.Select()
		return nil
	})
}
