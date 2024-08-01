package utils

import (
	"errors"
	"math"
	"reflect"
	v8 "rogchap.com/v8go"
	"testing"
)

func TestBaseTypeConvertToJsValue(t *testing.T) {
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
	}
	ctx := v8.NewContext()
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
			if !reflect.DeepEqual(gRes, tt.want) {
				t.Errorf("ToJsValue() gotRes(%v) = %v, want(%v) %v",
					reflect.TypeOf(gRes).String(), gRes,
					reflect.TypeOf(tt.want).String(), tt.want)
			}
		})
	}
}
