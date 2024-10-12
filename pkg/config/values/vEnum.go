package values

import "bytes"

type VEnum struct {
	value IValue
	// Enums *[]string
}

func (ve *VEnum) String() string {
	return ve.value.String()
}

func (ve *VEnum) ValueFromString(s string) (err error) {
	return ve.ValueFromString(s)
}

func (ve *VEnum) SchemaType() VType {
	return "enum" + "|" + ve.value.SchemaType()
}

func (ve *VEnum) Encode(buf bytes.Buffer) error {
	return ve.value.Encode(buf)
}

func (ve *VEnum) Decode(buf bytes.Buffer) error {
	return ve.value.Decode(buf)
}
