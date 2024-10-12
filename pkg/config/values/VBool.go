package values

import (
	"strconv"
)

type VBool struct {
	value bool
}

func (vb *VBool) String() string {
	return strconv.FormatBool(vb.value)
}

func (vb *VBool) ValueFromString(s string) (err error) {
	vb.value, err = strconv.ParseBool(s)
	return
}

func (vb *VBool) SchemaType() VType {
	return VBooleanT
}

func (vb *VBool) Value() any {
	return vb.value
}

func (vb *VBool) SetValue(val any) {
	vb.value = val.(bool)
}
