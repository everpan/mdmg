package values

import (
	"bytes"
	"fmt"
	"strings"
)

type IValue interface {
	String() string
	ValueFromString(string) error
	SchemaType() string
	Encode(buf bytes.Buffer) error
	Decode(buf bytes.Buffer) error
}

// CreateValue value factory
func CreateValue(typ string, strVal string) (IValue, error) {
	var val IValue
	if strings.Index(typ, "|") > -1 {
		return createCompositeValue(typ, strVal)
	}
	switch typ {
	case "boolean", "bool":
		val = &VBool{}
	case "float":
		val = &VFloat{}
	case "string":
		val = &VString{}
	case "number":
		val = &VNumber{}
	case "array", "enum":
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

func createCompositeValue(typ string, strVal string) (IValue, error) {
	var (
		val IValue
		err error
	)
	split := strings.Split(typ, "|")
	// check size ?
	typ1, typ2 := split[0], split[1]
	if typ1 == "enum" {
		v, err := CreateValue(typ2, strVal)
		if err != nil {
			return nil, err
		}
		val = &VEnum{
			value: v,
		}
	} else if typ1 == "array" {
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
