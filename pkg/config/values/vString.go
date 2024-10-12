package values

import (
	"bytes"
	"encoding/json"
)

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

func (vs *VString) Encode(buf bytes.Buffer) error {
	err := json.NewEncoder(&buf).Encode(vs.value)
	buf.Truncate(buf.Len() - 1)
	return err
}

func (vs *VString) Decode(buf bytes.Buffer) error {
	err := json.NewDecoder(&buf).Decode(&vs.value)
	return err
}
