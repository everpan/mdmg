package ctx

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
)

type ICodeResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func NewICodeResponse(code int, message string, data interface{}) *ICodeResponse {
	return &ICodeResponse{Code: code, Message: message, Data: data}
}

func (resp *ICodeResponse) Marshal() (data []byte) {
	data, _ = json.Marshal(resp)
	return
}

func (resp *ICodeResponse) Unmarshal(data []byte) {
	_ = json.Unmarshal(data, resp)
}

func SendInternalServerError(fc *fiber.Ctx, err error) error {
	fc.Response().Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	return SendError(fc, fiber.StatusInternalServerError, err)
}

func SendError(fc *fiber.Ctx, status int, e error) error {
	_ = fc.SendStatus(status)
	resp := NewICodeResponse(-1, e.Error(), nil)
	return fc.Send(resp.Marshal())
}

func SendSuccess(fc *fiber.Ctx, data any) error {
	resp := NewICodeResponse(0, "", data)
	return fc.Send(resp.Marshal())
}

func AppRouterAdd(router fiber.Router, h *IcPathHandler) {
	if h.Method != "" {
		router.Add(h.Method, h.Path, h.Handler.WrapHandler())
	} else {
		router.Group(h.Path, h.Handler.WrapHandler())
	}
}

func AppRouterAddMulti(router fiber.Router, handlers []*IcPathHandler) {
	for _, handler := range handlers {
		AppRouterAdd(router, handler)
	}
}

func AppRouterAddGroup(app *fiber.App, g *IcGroupPathHandler) {
	r := app.Group(g.GroupPath)
	AppRouterAddMulti(r, g.Handlers)
}
