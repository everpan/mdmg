package v8runtime

import (
	"github.com/everpan/mdmg/utils"
	"github.com/gofiber/fiber/v2"
	v8 "rogchap.com/v8go"
)

const (
	OutBaseURL int = iota
	OutOriginURL
	OutHeader
	OutParam
	OutQuery
)

func httpValue(fb *fiber.Ctx, iso *v8.Isolate, typ int) *v8.FunctionTemplate {
	return v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		if len(info.Args()) < 1 {
			return getAllValues(info, typ, fb)
		}
		return getValueByKey(info, typ, fb)
	})
}

func getAllValues(info *v8.FunctionCallbackInfo, typ int, fb *fiber.Ctx) *v8.Value {
	var rv *v8.Value
	var err error
	if OutHeader == typ {
		hs := fb.GetReqHeaders()
		rv, err = utils.ToJsValue(info.Context(), hs)
	} else if OutParam == typ {
		mm := fb.AllParams()
		rv, err = utils.ToJsValue(info.Context(), mm)
	} else if OutQuery == typ {
		mm := fb.Queries()
		rv, err = utils.ToJsValue(info.Context(), mm)
	}
	if err != nil {
		return utils.JsException(info.Context(), err)
	}
	return rv
}

func getValueByKey(info *v8.FunctionCallbackInfo, typ int, fb *fiber.Ctx) *v8.Value {
	var rv *v8.Value
	var err error
	k := info.Args()[0].String()
	var v string
	if OutHeader == typ {
		v = fb.Get(k)
	} else if OutParam == typ {
		v = fb.Params(k)
	} else if OutQuery == typ {
		v = fb.Query(k)
	}
	rv, err = v8.NewValue(info.Context().Isolate(), v)
	if err != nil {
		return utils.JsException(info.Context(), err)
	}
	return rv
}

func next(fb *fiber.Ctx, iso *v8.Isolate) *v8.FunctionTemplate {
	return v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		err := fb.Next()
		if err != nil {
			return utils.JsException(info.Context(), err)
		}
		jv, _ := v8.NewValue(iso, true)
		return jv
	})
}

func baseURL(fb *fiber.Ctx, iso *v8.Isolate, typ int) *v8.FunctionTemplate {
	return v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		var v string
		if OutBaseURL == typ {
			v = fb.BaseURL()
		} else if OutOriginURL == typ {
			v = fb.OriginalURL()
		}
		jv, err := v8.NewValue(iso, v)
		if err != nil {
			return utils.JsException(info.Context(), err)
		}
		return jv
	})
}

func ModuleVersion(c *Ctx, iso *v8.Isolate, f int) *v8.FunctionTemplate {
	return v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) (t *v8.Value) {
		m, v := utils.SplitModuleVersion(c.ModuleVersion)
		if f == 0 {
			t, _ = v8.NewValue(iso, m)
		} else {
			t, _ = v8.NewValue(iso, v)
		}
		return
	})
}

func ExportObject(c *Ctx, iso *v8.Isolate) *v8.ObjectTemplate {
	fb := c.fbCtx
	t := v8.NewObjectTemplate(iso)
	_ = t.Set("header", httpValue(fb, iso, OutHeader))
	// t.Set("next", next(fb, iso))
	_ = t.Set("query", httpValue(fb, iso, OutQuery))
	_ = t.Set("param", httpValue(fb, iso, OutParam))
	_ = t.Set("baseURL", baseURL(fb, iso, OutBaseURL))
	_ = t.Set("originURL", baseURL(fb, iso, OutOriginURL))
	// module() version() fetch module and version
	_ = t.Set("module", ModuleVersion(c, iso, 0))
	_ = t.Set("version", ModuleVersion(c, iso, 1))
	//
	jv, _ := v8.NewValue(iso, fb.BaseURL())
	_ = t.Set("baseURL2", jv)
	return t
}
