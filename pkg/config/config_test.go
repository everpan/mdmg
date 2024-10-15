package config

import (
	"encoding/json"
	. "github.com/everpan/mdmg/pkg/config/values"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOutputSchema(t *testing.T) {
	c := &IcSectionConfig{}
	schema := c.ExportSchema()
	c2 := &IcSectionConfig{}
	e := c2.ImportSchema(schema)
	assert.Nil(t, e)
	schema2 := c2.ExportSchema()
	assert.Equal(t, string(schema), string(schema2))
}

func TestParseSchema(t *testing.T) {
	secConf := GlobalConfig.NewSection("test", "test")
	secConf.AddEnumSchema(VBooleanT, "enable", "是否启用", "true", EnumDesc{"true": "启用", "false": "停用"})
	defConfig := string(secConf.ExportSchema())
	// t.Log(defConfig)
	enableVal := secConf.GetValue("enable")
	assert.Equal(t, true, enableVal)
	tests := []struct {
		name           string
		jsonConfig     string
		wantJsonConfig string
		wantErrString  string
	}{
		{"empty and json unmarshal error", "", `{
 "__desc": "empty and json unmarshal error"
}`, ""},
		{"empty and no error", "{}", `{
 "__desc": "empty and no error"
}`, ""},
		{"default config", defConfig, defConfig, ""}, // 由于map无序的参与，测试不稳定；注意甄别
	}
	conf := NewIConfig()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := conf.NewSection(tt.name, tt.name)
			e := c.ImportSchema([]byte(tt.jsonConfig))
			if tt.wantErrString == "" {
				assert.Nil(t, e)
			} else {
				assert.Contains(t, e.Error(), tt.wantErrString)
			}
			assert.Equal(t, tt.wantJsonConfig, string(c.ExportSchema()))
			// t.Log(string(conf.ExportValues()))
		})
	}
}

func TestEnumSchema(t *testing.T) {
	conf := NewIConfig()
	testEnum := NewItemSchema(CompositeT(VEnumT, VBoolT), "test", "just for test", "", "true")
	jd, _ := json.Marshal(testEnum)
	t.Log(string(jd), testEnum.Default)
	assert.Nil(t, testEnum.Value)
	assert.NotNil(t, testEnum.Default)
	jd, _ = json.Marshal(testEnum.Default)
	t.Log("enum value is true ", string(jd), testEnum.Default)
	// string
	err := conf.
		NewSection("test", "test").
		AddEnumSchema(VStringT, "str-enum", "str-enum-test", "g",
			EnumDesc{"a": "values of a", "b": "values of b"})
	assert.Contains(t, err.Error(), "must be one of [")

}

func TestConfigSchema(t *testing.T) {
	conf := NewIConfig()
	sysSec := conf.NewSection("system", "some config of system")
	sysSec.AddBooleanSchema("using-feature", "使用某种特性", "true")
	sysSec.AddBooleanSchema("using-feature-default-false ", "使用某种特性2", "false")

	sc := sysSec.GetSchema("using-feature")
	assert.NotNil(t, sc)
	assert.Equal(t, true, sc.Default)
	// assert.False(t, sc.IsSetVal)

	schemaJson := conf.ExportValues()
	t.Log(string(schemaJson))
	schemaJson2 := conf.ExportValuesFlat()
	assert.Contains(t, string(schemaJson2), "system.using-feature")

	assert.Equal(t, true, conf.GetValue("system.using-feature"), "")
	assert.Equal(t, false, conf.GetValue(" system.using-feature-default-false"), "")
	// set value

	assert.Nil(t, sc.Value)
	err := conf.SetValue("system.using-feature-default-false", true)
	assert.Nil(t, err)
	assert.Equal(t, true, conf.GetValue(" system.using-feature-default-false"), "")
	conf.NewSection("i-am-a-new-section", "new section")
	t.Log(string(conf.ExportSchema()))
}
