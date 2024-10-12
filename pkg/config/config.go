package config

import (
	"bytes"
	"encoding/json"
	"github.com/everpan/mdmg/pkg/config/values"
)

type Config struct {
	JSModuleRootPath string
	// module-version/${JSModuleBeckEndDir}/script.js
	JSModuleBeckEndDir string
	HiddenModelVersion bool `yaml:"hidden-model-version"`
}

type ItemDesc struct {
	Item string `json:"-"` // the key
	Desc string `json:"desc"`
}

type EnumDesc map[string]string

/*
	func (itemDesc *ItemDesc) Encode() []byte {
		buf := new(bytes.Buffer)
		encoder := json.NewEncoder(buf)
		encoder.Encode(itemDesc.Item)
		buf.Truncate(buf.Len() - 1)
		buf.WriteString(": ")
		encoder.Encode(itemDesc.Desc)
		buf.Truncate(buf.Len() - 1)
		return buf.Bytes()
	}
*/
const VSectionT values.VType = "section"

type Schema struct {
	Item     string        `json:"-"` // the key
	Desc     string        `json:"desc"`
	Type     values.VType  `json:"type"`                // category string number enum
	EnumDesc EnumDesc      `json:"enum-desc,omitempty"` // if type is enum , give each desc of enum value, item is the enum value.
	Value    values.IValue `json:"value,omitempty"`
	Default  values.IValue `json:"default,omitempty"`
	IsSetVal bool          `json:"-"`
}

type SchemaMap = map[string]*Schema

type IcConfig struct {
	section string
	sMap    SchemaMap
}

var (
	ICodeGlobalConfig = NewIConfig()
)

func (sc *Schema) GetValue() any {
	if sc.Value == nil {
		if sc.Default != nil {
			return sc.Default.Value()
		}
	} else {
		return sc.Value.Value()
	}
	return nil
}

func (sc *Schema) Update(schema *Schema, forceUpdateValue bool) {
	sc.Desc = schema.Desc
	sc.Type = schema.Type
	sc.EnumDesc = schema.EnumDesc
	if forceUpdateValue {
		sc.Value = schema.Value
	}
	sc.Default = schema.Default
}

func NewIConfig() *IcConfig {
	return &IcConfig{
		sMap: SchemaMap{},
	}
}
func NewItemSchema(section, item, desc string, typ values.VType, val, defVal string) *Schema {
	if typ == VSectionT {
		item = section
	} else if len(section) > 0 {
		item = section + "." + item
	}
	var (
		// err error
		sVal, dVal values.IValue
	)
	if val != "" {
		sVal, _ = values.CreateValue(typ, val)
	}
	dVal, _ = values.CreateValue(typ, defVal)

	return &Schema{
		item, desc, typ, nil, sVal, dVal, val != "",
	}
}

func init() {
	ICodeGlobalConfig.AddSection("system", "系统")
}
func (c *IcConfig) GetSchema(key string) *Schema {
	sc, ok := c.sMap[key]
	if !ok {
		return nil
	}
	return sc
}
func (c *IcConfig) GetValue(key string) (any, bool) {
	conf, ok := c.sMap[key]
	if ok {
		return conf.GetValue(), ok
	}
	return "", false
}

// OutputSchema 输出schema信息
func (c *IcConfig) OutputSchema() []byte {
	if len(c.sMap) == 0 {
		return []byte("{}")
	}
	buf := bytes.NewBufferString("{\n")
	for k, c := range c.sMap {
		buf.WriteString(" \"")
		buf.WriteString(k)
		buf.WriteString("\": ")
		data, _ := json.Marshal(c)
		buf.Write(data)
		buf.WriteString(",\n")
	}
	if len(c.sMap) > 0 {
		buf.Truncate(buf.Len() - 2)
	}
	buf.WriteString("\n}")
	return buf.Bytes()
}

// ParseSchema 解析schema，形成默认配置
func (c *IcConfig) ParseSchema(data []byte) error {
	return json.Unmarshal(data, &c.sMap)
}

func (c *IcConfig) LoadKeyValue(data []byte) error {
	return nil
}

func (c *IcConfig) OutputMap() map[string]any {
	kv := make(map[string]any)
	for _, item := range c.sMap {
		if item.Type == VSectionT {
			continue
		}
		kv[item.Item] = item.GetValue()
	}
	return kv
}

// OutputKeyValue 以kv的形式输出配置项; section is ignore
func (c *IcConfig) OutputKeyValue() []byte {
	kv := c.OutputMap()
	data, _ := json.MarshalIndent(kv, "", " ")
	return data
}

func (c *IcConfig) AddSchema(schema *Schema) {
	c.sMap[schema.Item] = schema
}
func (c *IcConfig) ReplaceSchema(schema *Schema) {
	sc, ok := c.sMap[schema.Item]
	if !ok {
		c.AddSchema(schema)
		return
	}
	sc.Update(schema, false)
}
func (c *IcConfig) ReplaceSchemaAndValue(schema *Schema) {
	sc, ok := c.sMap[schema.Item]
	if !ok {
		c.AddSchema(schema)
		return
	}
	sc.Update(schema, true)
}
func (c *IcConfig) AddSection(sec, desc string) *IcConfig {
	c.ReplaceSchema(NewItemSchema(sec, "", desc, VSectionT, "", ""))
	return &IcConfig{
		section: sec,
		sMap:    c.sMap,
	}
}
func (c *IcConfig) AddStringSchema(item, desc, defaultValue string) {
	c.ReplaceSchema(NewItemSchema(c.section, item, desc, values.VStringT, "", defaultValue))
}
func (c *IcConfig) AddNumberSchema(item, desc, defaultValue string) {
	c.ReplaceSchema(NewItemSchema(c.section, item, desc, values.VNumberT, "", defaultValue))
}
func (c *IcConfig) AddFloatSchema(item, desc, defaultValue string) {
	c.ReplaceSchema(NewItemSchema(c.section, item, desc, values.VFloatT, "", defaultValue))
}
func (c *IcConfig) AddBooleanSchema(item, desc, defaultValue string) {
	c.ReplaceSchema(NewItemSchema(c.section, item, desc, values.VBooleanT, "false", defaultValue))
}

func (c *IcConfig) AddEnumSchema(item, desc string, valType values.VType, defaultValue string, enumDesc EnumDesc) {
	schema := NewItemSchema(c.section, item, desc, values.CompositeT(values.VEnumT, valType), "", defaultValue)
	schema.EnumDesc = enumDesc
	c.ReplaceSchema(schema)
}

var DefaultConfig = Config{
	JSModuleRootPath:   "web/script_module",
	JSModuleBeckEndDir: "backend",
}
