package config

import (
	"bytes"
	"encoding/json"
	"strconv"
)

type IValue interface {
	String() string
	ValueFromString(string) error
	SchemaType() string
	Encode(buf bytes.Buffer) error
}

// VString string wrap
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

func (vs *VString) SchemaType() string {
	return "string"
}

func (vs *VString) Encode(buf bytes.Buffer) error {
	err := json.NewEncoder(&buf).Encode(vs.value)
	buf.Truncate(buf.Len() - 1)
	return err
}

// VNumber number wrap
type VNumber struct {
	value int64
}

func (vn *VNumber) String() string {
	return strconv.FormatInt(vn.value, 10)
}

func (vn *VNumber) ValueFromString(s string) (err error) {
	vn.value, err = strconv.ParseInt(s, 10, 64)
	return
}

func (vn *VNumber) SchemaType() string {
	return "number"
}

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

func (vf *VFloat) SchemaType() string {
	return "float"
}

type VArray struct {
	value []IValue
}

func (av *VArray) String() string {

}

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

func (ve *VEnum) SchemaType() string {
	return "enum"
}
