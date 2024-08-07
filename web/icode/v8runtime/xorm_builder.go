package v8runtime

import (
	"github.com/everpan/mdmg/utils"
	v8 "rogchap.com/v8go"
	"sync"
	"sync/atomic"
	"xorm.io/builder"
)

var builderPool = sync.Map{}
var builderKey uint32

const (
	OpGroupByString int = iota
	OpOrderByAny
	OpOrderByAsc
	OpHavingAny
)

func ExportXormBuilder(iso *v8.Isolate) *v8.ObjectTemplate {
	obj := v8.NewObjectTemplate(iso)
	obj.Set("select", select_(iso))
	obj.Set("orderBy", orderBy(iso))
	return obj
}

func select_(iso *v8.Isolate) *v8.FunctionTemplate {
	return v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		// info.This()
		argsStr := make([]string, len(info.Args()))
		for i, arg := range info.Args() {
			argsStr[i] = arg.String()
		}
		b := builder.Select(argsStr...)
		key := atomic.AddUint32(&builderKey, 1)
		obj := ExportXormBuilder(iso)
		storeBuilder(key, b)
		obj.Set("__builder_instance", key)
		r, err := v8.NewValue(iso, obj)
		if err != nil {
			return utils.JsException(info.Context(), err.Error())
		}
		return r
	})
}

func orderBy(iso *v8.Isolate) *v8.FunctionTemplate {
	return v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		builder, value, done := AcquireBuilderFromObject(info)
		if done {
			return value
		}
		if nil == builder {
			return utils.JsException(info.Context(), "sql builder is nil")
		}
		args, _ := utils.ToGoValues(info.Context(), info.Args())
		builder.OrderBy(args)
		r, _ := v8.NewValue(iso, info.This())
		return r
	})
}

func AcquireBuilderFromObject(info *v8.FunctionCallbackInfo) (*builder.Builder, *v8.Value, bool) {
	k, err := info.This().Get("__builder_instance")
	if err != nil {
		return nil, utils.JsException(info.Context(), err.Error()), true
	}
	builder := acquireBuilder(k.Uint32())
	return builder, nil, false
}

func acquireBuilder(k uint32) *builder.Builder {
	b, ok := builderPool.Load(k)
	if !ok {
		return nil
	}
	return b.(*builder.Builder) //  &builder.Builder{}
}

func storeBuilder(k uint32, builder *builder.Builder) {
	builderPool.Store(k, builder)
}
