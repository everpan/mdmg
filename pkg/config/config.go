package config

import (
	"encoding/json"
	"fmt"
	. "github.com/everpan/mdmg/pkg/config/values"
	"strings"
)

type ItemDesc struct {
	Item string `json:"-"` // the key
	Desc string `json:"desc"`
}
type EnumDesc map[string]string

type Schema struct {
	Item     string   `json:"-"` // the key
	Desc     string   `json:"desc"`
	Type     VType    `json:"type"`                // category string number enum
	EnumDesc EnumDesc `json:"enum-desc,omitempty"` // if type is enum , give each desc of enum value, item is the enum value.
	Value    IValue   `json:"-"`
	// Default  IValue `json:"default,omitempty"`
	Default any `json:"default,omitempty"`
}
type schemaExportValue struct {
	Desc     string            `json:"desc"`
	Type     string            `json:"type"`
	EnumDesc map[string]string `json:"enum-desc,omitempty"`
	Value    any               `json:"value,omitempty"`
	Default  any               `json:"default,omitempty"`
}

func (sc *Schema) toSchemaExportValue() *schemaExportValue {
	sev := &schemaExportValue{
		Desc:     sc.Desc,
		Type:     string(sc.Type),
		EnumDesc: sc.EnumDesc,
		Default:  sc.Default,
	}
	if sc.Value != nil {
		sev.Value = sc.Value.Value()
	}
	return sev
}

func (sc *Schema) fromSchemaExportValue(sev *schemaExportValue) {
	sc.Desc = sev.Desc
	sc.Type = VType(sev.Type)
	sc.EnumDesc = sev.EnumDesc
	sc.Default = sev.Default
	sc.Value, _ = CreateValue(sc.Type, fmt.Sprintf("%v", sev.Value))
}

func (sc *Schema) toJson() []byte {
	sev := sc.toSchemaExportValue()
	jd, _ := json.Marshal(sev)
	return jd
}

func (sc *Schema) GetValue() any {
	if sc.Value != nil {
		return sc.Value.Value()
	}
	return sc.Default
}

func (sc *Schema) Update(schema *Schema, updateValue bool) {
	sc.Desc = schema.Desc
	sc.Type = schema.Type
	sc.EnumDesc = schema.EnumDesc
	if updateValue {
		sc.Value = schema.Value
	}
	sc.Default = schema.Default
}

func (sc *Schema) SetValue(val any) error {
	v, e := CreateValue(sc.Type, fmt.Sprintf("%v", val))
	if e != nil {
		return e
	}
	sc.Value = v
	return nil
}

func (sc *Schema) ClearValue() {
	sc.Value = nil
}

func NewIConfig() *IcConfig {
	return &IcConfig{
		confMap: make(map[string]*IcSectionConfig),
	}
}

func NewItemSchema(typ VType, item, desc, val, defVal string) *Schema {
	var (
		// err error
		sVal, dVal IValue
		dv         any
	)
	if val != "" {
		sVal, _ = CreateValue(typ, val)
	}
	dVal, _ = CreateValue(typ, defVal)
	if dVal != nil {
		dv = dVal.Value()
	}
	return &Schema{strings.TrimSpace(item),
		desc, typ, nil, sVal, dv}
}

func init() {
	GlobalConfig.NewSection("system", "系统")
}

func (c *IcConfig) NewSection(section, desc string) *IcSectionConfig {
	seConf, ok := c.confMap[section]
	if !ok {
		seConf = &IcSectionConfig{description: desc, sMap: make(map[string]*Schema)}
		c.confMap[section] = seConf
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
		return fmt.Errorf("no schema of section: %s", sec)
	}
	return sc.SetValue(subKey, val)
}

func (c *IcConfig) ExportValuesMap() any {
	m := make(map[string]any)
	for k, v := range c.confMap {
		m[k] = v.ExportKeyValueMap()
	}
	return m
}

func (c *IcConfig) ExportValues() []byte {
	m := c.ExportValuesMap()
	d, _ := json.MarshalIndent(m, "", "  ")
	return d
}

// ExportValuesFlat 将所有的配置项扁平化为 "key"=value 的形式
func (c *IcConfig) ExportValuesFlat() []byte {
	m := make(map[string]any)
	for k, v := range c.confMap {
		m2 := v.ExportKeyValueMap()
		for k2, v2 := range m2 {
			m[k+"."+k2] = v2
		}
	}
	d, _ := json.MarshalIndent(m, "", "  ")
	return d
}

func (c *IcConfig) ExportSchema() []byte {
	var m = make(map[string]json.RawMessage)
	m["__config"] = json.RawMessage(`{"version": "1.0.0","format": "json"}`)
	for k, v := range c.confMap {
		m[k] = v.ExportSchema()
	}
	d, _ := json.MarshalIndent(m, "", "  ")
	return d
}

func (c *IcConfig) ImportSchema(data []byte) error {
	var m = make(map[string]json.RawMessage)
	err := json.Unmarshal(data, &m)
	if err != nil {
		return err
	}
	// todo parse by format, default is json
	for k, v := range m {
		sf := &IcSectionConfig{}
		err = sf.ImportSchema(v)
		c.confMap[k] = sf
		if err != nil {
			return err
		}
	}
	return nil
}
