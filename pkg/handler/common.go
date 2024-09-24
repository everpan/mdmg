package handler

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
	json.Unmarshal(data, resp)
}

func SendInternalServerError(fc *fiber.Ctx, err error) error {
	fc.Response().Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	return SendError(fc, fiber.StatusInternalServerError, err)
}

func SendError(fc *fiber.Ctx, status int, e error) error {
	fc.SendStatus(status)
	resp := NewICodeResponse(-1, e.Error(), nil)
	return fc.Send(resp.Marshal())
}
