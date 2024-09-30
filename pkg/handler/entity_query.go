package handler

import (
	"github.com/everpan/mdmg/pkg/ctx"
)

var EntityQueryHandler = &ctx.IcGroupPathHandler{
	GroupPath: "/entity",
	Handlers:  make([]*ctx.IcPathHandler, 0),
}

func query(c *ctx.IcContext) error {
	return nil
}
