package zcode

type FlowConfig struct {
	Input  string `json:"input"`
	Output string `json:"output"`
	From   string `json:"from"`
	To     string `json:"to"`
	Trace  string `json:"trace"`
	Next   bool   `json:"next"`
}
