package values

import (
	"bytes"
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

func (vn *VNumber) SchemaType() string {
	return "number"
}

func (vn *VNumber) Encode(buf bytes.Buffer) error {
	buf.WriteString(strconv.FormatInt(vn.value, 10))
	buf.WriteString(",\n")
	return nil
}

func (vn *VNumber) Decode(buf bytes.Buffer) error {
	line, err := buf.ReadString('\n')
	if err != nil {
		return err
	}
	vn.ValueFromString(line)
	return nil
}
