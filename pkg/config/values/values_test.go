package values

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIValueJson(t *testing.T) {
	var v IValue
	assert.Nil(t, v)
	d, e := json.Marshal(v)
	assert.Nil(t, e)
	assert.Equal(t, string(d), `null`)
}

func TestCreateValue(t *testing.T) {
	tests := []struct {
		name   string
		typ    string
		strVal string
		//wantTyp string
		wantVal string
		errStr  string
	}{
		{"not defined", "not defined", "", "", "invalid value type: not defined"},
		{"string", "string", "string - val - \"", "string - val - \"", ""},
		{"not number", "number", "not number", "", "invalid syntax"},
		{"number -1", "number", "-1", "-1", ""},
		{"number 0", "number", "0", "0", ""},
		{"number 9", "number", "9", "9", ""},
		{"number 9 with front space", "number", " 9", "9", ""},
		{"number 9 with back space", "number", "9 ", "9", ""},
		{"invalid number 9e", "number", "9e10 ", "9000", "invalid syntax"},
		{"bool 0", "boolean", "0", "false", ""},
		{"bool True", "bool", "True", "true", ""},
		{"bool False", "bool", "f", "false", ""},
		{"float False", "float", "4.3987", "4.3987", ""},
		{"float invalid", "float", "4.3987f", "4.3987", "invalid syntax"},
		{"array string", "array|string", "12,3,4,5", "12,3,4,5", ""},
		{"array bool", "array|boolean", "0,f,False,FALSE,1,t,True,TRUE",
			"false,false,false,false,true,true,true,true", ""},
		{"just array", "array", "", "", "using sample: array|number"},
		{"array|float", "array|float", "4.1,5.7", "4.1,5.7", ""},
		{"array|number", "array|number", "4,5,7", "4,5,7", ""},
		{"enum|number", "enum|number", "489", "489", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := CreateValue(VType(tt.typ), tt.strVal)
			if err != nil {
				assert.Nil(t, val)
				assert.NotEmptyf(t, tt.errStr, "return error: %v", err) // 错误信息必须给定
				assert.Contains(t, err.Error(), tt.errStr)
				return
			}
			assert.NotNil(t, val)
			if tt.typ == "bool" {
				tt.typ = "boolean"
			}
			assert.Equalf(t, tt.typ, string(val.SchemaType()), "type %v, %v", tt.typ, val.SchemaType())
			// t.Log(tt.wantVal, val.String())
			assert.Equalf(t, tt.wantVal, val.String(), "CreateValue(%v, %v)", tt.typ, tt.wantVal)
		})
	}
}
