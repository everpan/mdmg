package values

import (
	"strconv"
	"strings"
)

type VNumber struct {
	value int64
}

func (vn *VNumber) String() string {
	return strconv.FormatInt(vn.value, 10)
}

func (vn *VNumber) ValueFromString(s string) (err error) {
	vn.value, err = strconv.ParseInt(
		strings.TrimSpace(s), 10, 64)
	return
}

func (vn *VNumber) SchemaType() VType {
	return VNumberT
}

func (vn *VNumber) Value() any {
	return vn.value
}
