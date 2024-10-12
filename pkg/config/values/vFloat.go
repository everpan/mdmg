package values

import (
	"bytes"
	"strconv"
)

type VFloat struct {
	value float64
}

func (vf *VFloat) String() string {
	return strconv.FormatFloat(vf.value, 'f', -1, 64)
}

func (vf *VFloat) ValueFromString(s string) (err error) {
	vf.value, err = strconv.ParseFloat(s, 64)
	return
}

func (vf *VFloat) SchemaType() VType {
	return VFloatT
}

func (vf *VFloat) Encode(buf bytes.Buffer) error {
	return nil
}
func (vf *VFloat) Decode(buf bytes.Buffer) error {
	return nil
}
