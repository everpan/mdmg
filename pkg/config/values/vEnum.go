package values

type VEnum struct {
	value IValue
	// Enums *[]string
}

func (ve *VEnum) String() string {
	return ve.value.String()
}

func (ve *VEnum) ValueFromString(s string) (err error) {
	return ve.value.ValueFromString(s)
}

func (ve *VEnum) SchemaType() VType {
	return "enum" + "|" + ve.value.SchemaType()
}

func (ve *VEnum) Value() any {
	return ve.value.Value()
}

//func (ve *VEnum) MarshalJSON() ([]byte, error) {
//	return []byte(`{"everpan":true}`), nil
//}
//func (ve *VEnum) UnmarshalJSON(data []byte) error {
//	return nil
//}
