package utils

import (
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	v8 "rogchap.com/v8go"
)

type FuncT struct {
	c *v8.Context
	v *v8.Value
}

type PromiseT struct {
	c *v8.Context
	v *v8.Value
}

type UndefinedT byte

var Undefined = UndefinedT(0)

func ToJsValues(c *v8.Context, gVals []any) (res []*v8.Value, err error) {
	res = make([]*v8.Value, 0, len(gVals))
	for _, gVal := range gVals {
		jVal, err := ToJsValue(c, gVal)
		if err != nil {
			return nil, err
		}
		res = append(res, jVal)
	}
	return res, nil
}

func ToJsValue(c *v8.Context, gVal any) (res *v8.Value, err error) {
	iso := c.Isolate()
	if nil == gVal {
		return v8.Null(iso), nil
	}
	switch v := gVal.(type) {
	case string, bool, *big.Int, float64:
		return v8.NewValue(iso, v)

	case int8:
		return v8.NewValue(iso, int32(v))
	case int16:
		return v8.NewValue(iso, int32(v))
	case int32:
		return v8.NewValue(iso, v)
	case int:
		return v8.NewValue(iso, int32(v))
	case int64:
		return v8.NewValue(iso, v)
		// unsigned
	case uint8:
		return v8.NewValue(iso, int32(v))
	case uint16:
		return v8.NewValue(iso, int32(v))
	case uint32:
		return v8.NewValue(iso, int64(v))
	case uint64:
		if v <= math.MaxInt64 {
			return v8.NewValue(iso, int64(v))
		} else {
			return v8.NewValue(iso, float64(v))
		}
	case float32:
		return v8.NewValue(iso, float64(v))

	case FuncT:
		return gVal.(FuncT).v, nil
	case PromiseT:
		return gVal.(PromiseT).v, nil
	default:
		return ToJsValue(c, v)
	}
}

func ToGoValues(c *v8.Context, jVals []*v8.Value) (res []any, err error) {
	gVals := make([]any, 0, len(jVals))
	for _, jVal := range jVals {
		gVal, err := ToGoValue(c, jVal)
		if err != nil {
			return nil, err
		}
		gVals = append(gVals, gVal)
	}
	return gVals, nil
}

func ToGoValue(c *v8.Context, jVal *v8.Value) (any, error) {
	if nil == jVal || jVal.IsNull() {
		return nil, nil
	}
	if jVal.IsUndefined() {
		return Undefined, nil
	}
	if jVal.IsString() {
		return jVal.String(), nil
	}
	if jVal.IsBoolean() {
		return jVal.Boolean(), nil
	}
	if jVal.IsBigInt() {
		return jVal.BigInt().Int64(), nil
	}
	if jVal.IsNumber() {
		if jVal.IsInt32() {
			return int(jVal.Int32()), nil
		}
		return jVal.Number(), nil
	}
	//if jVal.IsSharedArrayBuffer() { // []byte
	//	buf, cleanFn, err := jVal.Sha
	//}
	if jVal.IsArray() {
		return goValueParse(jVal, []any{})
	}
	if jVal.IsMap() {
		return goValueParse(jVal, map[string]any{})
	}
	// Others
	var r any
	return goValueParse(jVal, r)
}

func goValueParse(jv *v8.Value, gv any) (any, error) {
	d, err := jv.MarshalJSON()
	if err != nil {
		return nil, err
	}
	p := &gv
	err = json.Unmarshal(d, p)
	if err != nil {
		return nil, err
	}
	return *p, err
}

func (fn *FuncT) String() string {
	return fmt.Sprintf("[Function: %s]", fn.v.String())
}

func (fn *FuncT) Call(args ...any) (any, error) {
	if fn.c == nil {
		return nil, fmt.Errorf("function context invalid")
	}
	cb, err := fn.v.AsFunction()
	if err != nil {
		return nil, err
	}
	jArgs, err := ToJsValues(fn.c, args)
	if err != nil {
		return nil, err
	}
	defer ReleaseJsValues(jArgs)
	val, err := cb.Call(fn.c.Global(), Valuer(jArgs)...)
	if err != nil {
		return nil, err
	}
	defer val.Release()
	gVal, err := ToGoValue(fn.c, val)
	if err != nil {
		return nil, err
	}
	return gVal, err
}

func Valuer(args []*v8.Value) []v8.Valuer {
	valuers := make([]v8.Valuer, len(args))
	for i, arg := range args {
		valuers[i] = arg
	}
	return valuers
}

func ReleaseJsValues(jVals []*v8.Value) {
	if jVals == nil || len(jVals) == 0 {
		return
	}
	for _, jVal := range jVals {
		if jVal.IsNull() || jVal.IsUndefined() {
			continue
		}
		jVal.Release()
	}
}

func (pro *PromiseT) String() string {
	p, e := pro.v.AsPromise()
	if e != nil {
		return "Promise: %s" + e.Error()
	}
	state := "pending"
	switch p.State() {
	case v8.Fulfilled:
		state = "fulfilled"
	case v8.Rejected:
		state = "rejected"
	}
	return fmt.Sprintf("[Promise: %s]", state)
}

func (u *UndefinedT) String() string {
	return "Undefined"
}
