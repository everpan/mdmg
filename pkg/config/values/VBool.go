package values

import (
	"bytes"
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

func (vb *VBool) SchemaType() string {
	return "boolean"
}

func (vb *VBool) Encode(buf bytes.Buffer) error {
	return nil
}

func (vb *VBool) Decode(buf bytes.Buffer) error {
	return nil
}
