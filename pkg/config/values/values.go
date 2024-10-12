package values

import (
	"bytes"
	"fmt"
	"strings"
)

type VType string

const (
	VArrayT   VType = "array"
	VBooleanT VType = "boolean"
	VBoolT    VType = "bool"
	VEnumT    VType = "enum"
	VFloatT   VType = "float"
	VNumberT  VType = "number"
	VStringT  VType = "string"
)

type IValue interface {
	String() string
	ValueFromString(string) error
	SchemaType() VType
	Encode(buf bytes.Buffer) error
	Decode(buf bytes.Buffer) error
}

// CreateValue value factory
func CreateValue(typ VType, strVal string) (IValue, error) {
	var val IValue
	if strings.Index(string(typ), "|") > -1 {
		return createCompositeValue(typ, strVal)
	}
	switch typ {
	case VBoolT, VBooleanT:
		val = &VBool{}
	case VFloatT:
		val = &VFloat{}
	case VStringT:
		val = &VString{}
	case VNumberT:
		val = &VNumber{}
	case VArrayT, VEnumT:
		return nil, fmt.Errorf("composite type %s, using sample: %s|number", typ, typ)
	default:
		return nil, fmt.Errorf("invalid value type: %s", typ)
	}
	err := val.ValueFromString(strVal)
	if err != nil {
		return nil, err
	}
	return val, err
}

func createCompositeValue(typ VType, strVal string) (IValue, error) {
	var (
		val IValue
		err error
	)
	split := strings.Split(string(typ), "|")
	// check size ?
	typ1, typ2 := VType(split[0]), VType(split[1])
	if typ1 == VEnumT {
		v, err := CreateValue(typ2, strVal)
		if err != nil {
			return nil, err
		}
		val = &VEnum{
			value: v,
		}
	} else if typ1 == VArrayT {
		val = &VArray{
			typ:   typ2,
			value: make([]IValue, 0),
		}
		err = val.ValueFromString(strVal)
	} else {
		return nil, fmt.Errorf("invalid value type: %s", typ1)
	}
	return val, err
}
