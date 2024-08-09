package handler

import (
	"github.com/everpan/mdmg/web/config"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type checkFun func(t *testing.T, r *http.Response, e error)

var echoBody = func(t *testing.T, r *http.Response, e error) {
	body, _ := io.ReadAll(r.Body)
	ret := &ICodeResponse{}
	ret.Unmarshal(body)
	t.Log("body", string(body))
	t.Log(ret)
}

var contains = func(code int, msg string) checkFun {
	return func(t *testing.T, r *http.Response, e error) {
		body, _ := io.ReadAll(r.Body)
		ret := &ICodeResponse{}
		ret.Unmarshal(body)
		// t.Log("body", string(body), ret)
		assert.Equal(t, code, ret.Code)
		assert.Contains(t, string(body), msg)
	}
}

var fileNotExist = func(t *testing.T, r *http.Response, e error) {
	assert.Equal(t, http.StatusNotFound, r.StatusCode)
	body, _ := io.ReadAll(r.Body)
	str := string(body)
	assert.Contains(t, str, "can not find ")
	assert.Contains(t, str, ".js")
}

var wantInternalServerError = func(msg string) func(*testing.T, *http.Response, error) {
	return func(t *testing.T, res *http.Response, err error) {
		if nil != err {
			t.Error(err)
		}
		body, _ := io.ReadAll(res.Body)
		ret := &ICodeResponse{}
		ret.Unmarshal(body)
		assert.Nilf(t, err, "body: %v", string(body))
		assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
		assert.Equal(t, msg, ret.Message)
	}
}

func TestIcodeHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		scriptFileName string
		check          checkFun
	}{
		{"file not exist", fiber.MethodGet, "not_exist",
			fileNotExist},
		{"undefined error", fiber.MethodGet, "ret_undefined",
			wantInternalServerError("v8go: value is not an Object")},
		{"not ret object", fiber.MethodGet, "not_object",
			wantInternalServerError("v8go: value is not an Object")},
		{"output_get", fiber.MethodGet, "output",
			contains(0, `"module":"test"`)},
		{"output_delete", fiber.MethodDelete, "output",
			contains(0, "delete is keyword in js, alias is del")},
		{"not found output", fiber.MethodPatch, "output",
			contains(-1, "output object is not found in response")},
		{"method not found", fiber.MethodPut, "output",
			contains(-1, "not found the handler of method(PUT)")},
		{"dir1/dir2/dir3", fiber.MethodPut, "sub1/sub2/output",
			contains(-1, "/sub1/sub2/output.js")},
	}
	app := fiber.New()
	// Path:    "/v1/icode/:modVer/:jsFile/*",
	app.Group(ICoderHandler.Path, ICoderHandler.Handler)
	config.DefaultConfig.JSModuleRootPath = "../script_module"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := "/v1/icode/test-0.1.0/" + strings.TrimSpace(tt.scriptFileName)
			req := httptest.NewRequest(tt.method, target, nil)
			req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
			resp, err := app.Test(req)
			if tt.check != nil {
				tt.check(t, resp, err)
			}
		})
	}
}
