package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/everpan/mdmg/pkg/config/values"
	"strings"
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

type Schema struct {
	Item     string        `json:"-"` // the key
	Desc     string        `json:"desc"`
	Type     values.VType  `json:"type"`                // category string number enum
	EnumDesc EnumDesc      `json:"enum-desc,omitempty"` // if type is enum , give each desc of enum value, item is the enum value.
	Value    values.IValue `json:"-"`
	// Default  values.IValue `json:"default,omitempty"`
	Default  any  `json:"default,omitempty"`
	IsSetVal bool `json:"-"`
}

type IcSectionConfig struct {
	description string
	sMap        map[string]*Schema
}

type IcConfig struct {
	confMap map[string]*IcSectionConfig
}

func (sc *Schema) GetValue() any {
	//if sc.Value == nil {
	//	if sc.Default != nil {
	//		return sc.Default
	//	}
	//} else {
	//	return sc.Value.Value()
	//}
	if sc.IsSetVal {
		return sc.Value.Value()
	}
	return sc.Default
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

func (sc *Schema) SetValue(val any) error {
	v, e := values.CreateValue(sc.Type, fmt.Sprintf("%v", val))
	if e != nil {
		return e
	}
	sc.Value = v
	return nil
}

func (sc *Schema) ClearValue() {
	sc.Value = nil
	sc.IsSetVal = false
}

func NewIConfig() *IcConfig {
	return &IcConfig{
		confMap: make(map[string]*IcSectionConfig),
	}
}

func NewItemSchema(item, desc string, typ values.VType, val, defVal string) *Schema {
	var (
		// err error
		sVal, dVal values.IValue
		dv         any
	)
	if val != "" {
		sVal, _ = values.CreateValue(typ, val)
	}
	dVal, _ = values.CreateValue(typ, defVal)
	if dVal != nil {
		dv = dVal.Value()
	}
	return &Schema{strings.TrimSpace(item),
		desc, typ, nil, sVal, dv, val != ""}
}

func init() {
	GlobalConfig.NewSection("system", "系统")
}

func (c *IcConfig) NewSection(sec, desc string) *IcSectionConfig {
	seConf, ok := c.confMap[sec]
	if !ok {
		seConf = &IcSectionConfig{description: desc, sMap: make(map[string]*Schema)}
		c.confMap[sec] = seConf
	}
	return seConf
}

func (c *IcConfig) Section(sec string) *IcSectionConfig {
	seConf, ok := c.confMap[sec]
	if !ok {
		return nil
	}
	return seConf
}

func SplitSectionKey(key string) (section, subKey string) {
	key1 := strings.TrimSpace(key)
	dotPos := strings.IndexByte(key1, '.')
	if dotPos > -1 {
		section = key1[:dotPos]
		subKey = key1[dotPos+1:]
		return
	}
	return key, ""
}

func (c *IcConfig) GetValue(key string) any {
	sec, subKey := SplitSectionKey(key)
	if subKey == "" {
		return nil
	}
	secConf := c.Section(sec)
	if secConf == nil {
		return nil
	}
	return secConf.GetValue(subKey)
}

func (c *IcConfig) SetValue(key string, val any) error {
	sec, subKey := SplitSectionKey(key)
	if subKey == "" {
		return fmt.Errorf("invalid section key: %s", key)
	}
	sc := c.Section(sec)
	if sc == nil {
		return fmt.Errorf("invalid schema of section: %s", sec)
	}
	return sc.SetValue(subKey, val)
}

func (c *IcConfig) OutputMap() any {
	m := make(map[string]any)
	for k, v := range c.confMap {
		m[k] = v.OutputMap()
	}
	return m
}

func (c *IcConfig) OutputValues() []byte {
	m := c.OutputMap()
	d, _ := json.MarshalIndent(m, "", "  ")
	return d
}

func (c *IcConfig) OutputFlatValues() []byte {
	m := make(map[string]any)
	for k, v := range c.confMap {
		m2 := v.OutputMap()
		for k2, v2 := range m2 {
			m[k+"."+k2] = v2
		}
	}
	d, _ := json.MarshalIndent(m, "", "  ")
	return d
}

func (c *IcSectionConfig) GetSchema(key string) *Schema {
	sc, ok := c.sMap[key]
	if !ok {
		return nil
	}
	return sc
}

func (c *IcSectionConfig) GetValue(key string) any {
	scheam, ok := c.sMap[key]
	if ok {
		return scheam.GetValue()
	}
	return nil
}

// OutputSchema 输出schema信息
func (c *IcSectionConfig) OutputSchema() []byte {
	if len(c.sMap) == 0 {
		return []byte("{}")
	}
	buf := bytes.NewBufferString("{\n")
	for k, schema := range c.sMap {
		buf.WriteString(" \"")
		buf.WriteString(k)
		buf.WriteString("\": ")
		data, _ := json.Marshal(schema)
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
func (c *IcSectionConfig) ParseSchema(data []byte) error {
	return json.Unmarshal(data, &c.sMap)
}

func (c *IcSectionConfig) LoadKeyValue(data []byte) error {
	return nil
}

func (c *IcSectionConfig) OutputMap() map[string]any {
	kv := make(map[string]any)
	for _, item := range c.sMap {
		kv[item.Item] = item.GetValue()
	}
	return kv
}

// OutputKeyValue 以kv的形式输出配置项; description is ignore
func (c *IcSectionConfig) OutputKeyValue() []byte {
	kv := c.OutputMap()
	data, _ := json.MarshalIndent(kv, "", " ")
	return data
}

func (c *IcSectionConfig) AddSchema(schema *Schema) {
	c.sMap[schema.Item] = schema
}
func (c *IcSectionConfig) ReplaceSchema(schema *Schema) {
	sc, ok := c.sMap[schema.Item]
	if !ok {
		c.AddSchema(schema)
		return
	}
	sc.Update(schema, false)
}
func (c *IcSectionConfig) ReplaceSchemaAndValue(schema *Schema) {
	sc, ok := c.sMap[schema.Item]
	if !ok {
		c.AddSchema(schema)
		return
	}
	sc.Update(schema, true)
}
func (c *IcSectionConfig) AddStringSchema(item, desc, defaultValue string) {
	c.ReplaceSchema(NewItemSchema(item, desc, values.VStringT, "", defaultValue))
}
func (c *IcSectionConfig) AddNumberSchema(item, desc, defaultValue string) {
	c.ReplaceSchema(NewItemSchema(item, desc, values.VNumberT, "", defaultValue))
}
func (c *IcSectionConfig) AddFloatSchema(item, desc, defaultValue string) {
	c.ReplaceSchema(NewItemSchema(item, desc, values.VFloatT, "", defaultValue))
}
func (c *IcSectionConfig) AddBooleanSchema(item, desc, defaultValue string) {
	c.ReplaceSchema(NewItemSchema(item, desc, values.VBooleanT, "", defaultValue))
}

func (c *IcSectionConfig) AddEnumSchema(item, desc string, valType values.VType, defaultValue string, enumDesc EnumDesc) {
	schema := NewItemSchema(item, desc, values.CompositeT(values.VEnumT, valType), "", defaultValue)
	schema.EnumDesc = enumDesc
	c.ReplaceSchema(schema)
}

func (c *IcSectionConfig) SetValue(key string, val any) error {
	sc := c.GetSchema(key)
	if sc == nil {
		return fmt.Errorf("no schema for key: %s", key)
	}
	sc.IsSetVal = true
	return sc.SetValue(val)
}

var DefaultConfig = Config{
	JSModuleRootPath:   "web/script_module",
	JSModuleBeckEndDir: "backend",
}
