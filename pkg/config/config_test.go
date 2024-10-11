package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOutputSchema(t *testing.T) {
	c := &IcConfig{}
	schema := c.OutputSchema()
	t.Log(schema)
	c2 := &IcConfig{}
	e := c2.ParseSchema(schema)
	assert.Nil(t, e)
	schema2 := c2.OutputSchema()
	assert.Equal(t, schema, schema2)
}

func TestParseSchema(t *testing.T) {
	ICodeGlobalConfig.AddEnumSchema("enable", "是否启用", "true",
		EnumDesc{"true": "启用", "false": "停用"})
	defConfig := string(ICodeGlobalConfig.OutputSchema())
	t.Log(defConfig)
	enableVal, _ := ICodeGlobalConfig.GetValue("enable")
	assert.Equal(t, "true", enableVal)
	tests := []struct {
		name           string
		jsonConfig     string
		wantJsonConfig string
		wantErrString  string
	}{
		{"empty and json unmarshal error", "", "{}", "unexpected end of JSON input"},
		{"empty and no error", "{}", "{}", ""},
		{"default config", defConfig, defConfig, ""}, // 由于map无序的参与，测试不稳定；注意甄别
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewIConfig()
			e := c.ParseSchema([]byte(tt.jsonConfig))
			if tt.wantErrString == "" {
				assert.Nil(t, e)
			} else {
				assert.Contains(t, e.Error(), tt.wantErrString)
			}
			assert.Equal(t, tt.wantJsonConfig, string(c.OutputSchema()))
		})
	}
}
