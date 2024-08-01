package utils

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	v8 "rogchap.com/v8go"
	"testing"
)

func TestBaseTypeConvertToJsValue(t *testing.T) {
	stru := struct {
		Name string `json:"name"`
		Age  int32  // converted to float64
	}{"golang", 15}

	tests := []struct {
		name    string
		gVal    any
		want    any
		wantErr error
	}{
		{"string", "string", "string", nil},
		{"int8", int8(-1), int(-1), nil},
		{"int16", int16(-1), int(-1), nil},
		{"int32", int32(-1), int(-1), nil},
		{"int", int(-1), int(-1), nil},
		{"int64", int64(-1), int64(-1), nil},
		{"int64 min", int64(-9223372036854775808), int64(-9223372036854775808), nil},
		// unsigned
		{"uint8", uint8(255), int(255), nil},
		{"uint16 max", uint16(65535), int(65535), nil},
		{"uint32 max", uint32(4294967295), int64(4294967295), nil},
		{"uint64 max", uint64(math.MaxInt64), int64(math.MaxInt64), nil},
		// 通过json方式转换之后，数字类型，最终变为了float64
		{"struct to map", stru, map[string]any{"Age": float64(stru.Age), "name": stru.Name}, nil},
		{"array string", []string{"a", "b", "c"}, []any{"a", "b", "c"}, nil},
		{"array mix", []any{"a", "b", "c", 123, 4.6}, []any{"a", "b", "c", 123.0, 4.6}, nil},
	}
	ctx := v8.NewContext()
	defer ctx.Isolate().Dispose()
	defer ctx.Close()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jRes, err := ToJsValue(ctx, tt.gVal)
			if err != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("ToJsValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gRes, err := ToGoValue(ctx, jRes)
			if err != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("ToJsValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// assert.Equal(t, gRes, tt.want)
			if !reflect.DeepEqual(gRes, tt.want) {
				t.Errorf("jVal string: %v\nToJsValue():\n gotRes(%v)\t= %v,\n want(%v)\t= %v",
					jRes.String(),
					reflect.TypeOf(gRes).String(), gRes,
					reflect.TypeOf(tt.want).String(), tt.want)
			}
		})
	}
}

func compareMap(m1, m2 map[string]any) bool {
	if len(m1) != len(m2) {
		return false
	}
	for k, v := range m1 {
		v2, ok := m2[k]
		fmt.Printf("%v %v == %v %v %v\n", k, v, v2, ok, v == v2)
		fmt.Printf("%v %v\n", reflect.TypeOf(v).String(), reflect.TypeOf(v2).String())
		if !ok || v != v2 {
			return false
		}
	}
	return true
}
