package values

import (
	"bytes"
	"strings"
)

type VArray struct {
	typ   VType
	value []IValue
}

func (av *VArray) String() string {
	var buf bytes.Buffer
	for _, v := range av.value {
		buf.WriteString(v.String())
		buf.WriteString(",")
	}
	buf.Truncate(buf.Len() - 1)
	return buf.String()
}

func (av *VArray) ValueFromString(s string) (err error) {
	spVal := strings.Split(s, ",")
	var val IValue
	for _, str := range spVal {
		val, err = CreateValue(av.typ, str)
		if err != nil {
			av.value = nil
			return err
		}
		av.value = append(av.value, val)
	}
	return nil
}
func (av *VArray) SchemaType() VType {
	if len(av.value) == 0 {
		return "array|null"
	}
	return "array" + "|" + av.value[0].SchemaType()
}
func (av *VArray) Encode(buf bytes.Buffer) error {
	return nil
}
func (av *VArray) Decode(buf bytes.Buffer) error {
	return nil
}

func (av *VArray) SetType(typ VType) {
	av.typ = typ
}
