package values

type VString struct {
	value string
}

func (vs *VString) String() string {
	return vs.value
}

func (vs *VString) ValueFromString(s string) error {
	vs.value = s
	return nil
}

func (vs *VString) SchemaType() VType {
	return VStringT
}

func (vs *VString) Value() any {
	return vs.value
}
