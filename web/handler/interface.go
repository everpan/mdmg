package handler

import (
	"github.com/everpan/mdmg/v8runtime"
	"os"
	"path/filepath"
	v8 "rogchap.com/v8go"
)

type RunScript interface {
	RunScript(ctx *v8runtime.Ctx, str string) (*v8.Value, error)
}
type runFileScript struct{}

func (r *runFileScript) RunScript(ctx *v8runtime.Ctx, str string) (*v8.Value, error) {
	script, err := os.ReadFile(str)
	if err != nil {
		return nil, err
	}
	str = filepath.Base(str)
	return ctx.RunScript(string(script), str)
}

type RunCode struct{}

func (r *RunCode) RunScript(ctx *v8runtime.Ctx, str string) (*v8.Value, error) {
	return ctx.RunScript(str, "_run_code.js")
}
