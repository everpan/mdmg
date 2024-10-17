package config

import (
	"encoding/json"
	"fmt"
	. "github.com/everpan/mdmg/pkg/config/values"
	"strings"
)

type IcSectionConfig struct {
	description string
	sMap        map[string]*Schema
}

type IcConfig struct {
	confMap map[string]*IcSectionConfig
}

func (c *IcSectionConfig) GetSchema(key string) *Schema {
	sc, ok := c.sMap[key]
	if !ok {
		return nil
	}
	return sc
}

func (c *IcSectionConfig) GetValue(key string) any {
	schema, ok := c.sMap[key]
	if ok {
		return schema.GetValue()
	}
	return nil
}

// ExportSchema 输出schema信息
func (c *IcSectionConfig) ExportSchema() []byte {
	var m = make(map[string]any)
	for k, v := range c.sMap {
		m[k] = v
	}
	m["__desc"] = c.description
	jd, _ := json.MarshalIndent(m, "", " ")
	return jd
}

// ImportSchema 解析schema，形成默认配置
func (c *IcSectionConfig) ImportSchema(data []byte) error {
	var m = make(map[string]json.RawMessage)
	json.Unmarshal(data, &m)
	json.Unmarshal(m["__desc"], &c.description)
	delete(m, "__desc")
	for k, v := range m {
		var schema = Schema{}
		json.Unmarshal(v, &schema)
		c.sMap[k] = &schema
	}
	return nil
}

func (c *IcSectionConfig) LoadKeyValue(data []byte) error {
	return nil
}

func (c *IcSectionConfig) ExportKeyValueMap() map[string]any {
	kv := make(map[string]any)
	for _, item := range c.sMap {
		kv[item.Item] = item.GetValue()
	}
	return kv
}

// ExportKeyValue 以kv的形式输出配置项; description is ignore
func (c *IcSectionConfig) ExportKeyValue() []byte {
	kv := c.ExportKeyValueMap()
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
	c.ReplaceSchema(NewItemSchema(VStringT, item, desc, "", defaultValue))
}
func (c *IcSectionConfig) AddNumberSchema(item, desc, defaultValue string) {
	c.ReplaceSchema(NewItemSchema(VNumberT, item, desc, "", defaultValue))
}
func (c *IcSectionConfig) AddFloatSchema(item, desc, defaultValue string) {
	c.ReplaceSchema(NewItemSchema(VFloatT, item, desc, "", defaultValue))
}
func (c *IcSectionConfig) AddBooleanSchema(item, desc, defaultValue string) {
	c.ReplaceSchema(NewItemSchema(VBooleanT, item, desc, "", defaultValue))
}

func (c *IcSectionConfig) SetDefault(item, desc string, value any) {

}
func (c *IcSectionConfig) SetEnumDefault(item, desc string, value any, enumDesc EnumDesc) {

}

func (c *IcSectionConfig) AddEnumSchema(valType VType, item, desc, defaultValue string, enumDesc EnumDesc) error {
	_, ok := enumDesc[defaultValue]
	if !ok {
		var allValue []string
		for ed := range enumDesc {
			allValue = append(allValue, ed)
		}
		return fmt.Errorf("default value `%s` must be one of [%v]",
			defaultValue, strings.Join(allValue, ","))
	}
	schema := NewItemSchema(CompositeT(VEnumT, valType), item, desc, "", defaultValue)
	schema.EnumDesc = enumDesc
	c.ReplaceSchema(schema)
	return nil
}

func (c *IcSectionConfig) SetValue(key string, val any) error {
	sc := c.GetSchema(key)
	if sc == nil {
		return fmt.Errorf("no schema for key: %s", key)
	}
	return sc.SetValue(val)
}
