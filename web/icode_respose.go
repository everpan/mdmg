package main

import "encoding/json"

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
