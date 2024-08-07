package v8runtime

import (
	"github.com/everpan/mdmg/utils"
	"github.com/everpan/mdmg/web/icode"
	v8 "rogchap.com/v8go"
	"xorm.io/xorm"
)

func ExportXormObject(ctx *icode.Ctx, iso *v8.Isolate) *v8.ObjectTemplate {
	obj := v8.NewObjectTemplate(iso)
	_ = obj.Set("exec", execSql(ctx, iso))
	_ = obj.Set("transaction_exec", transactionExec(ctx, iso))
	_ = obj.Set("tExec", transactionExec(ctx, iso))
	_ = obj.Set("query", queryInterface(ctx, iso))

	return obj
}

func execSql(ctx *icode.Ctx, iso *v8.Isolate) *v8.FunctionTemplate {
	return v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		c := info.Context()
		if len(info.Args()) < 0 {
			return utils.JsException(c, "no sql found")
		}
		args := JsArgsToGoArgs(info, c)
		eng := ctx.GetEngine()
		ret, err := eng.Exec(info.Args()[0].String(), args)
		if err != nil {
			return utils.JsError(c, err.Error())
		}
		var R struct {
			LastInsertId int64 `json:"last_insert_id"`
			RowsAffected int64 `json:"rows_affected"`
		}
		R.RowsAffected, _ = ret.LastInsertId()
		R.LastInsertId, _ = ret.RowsAffected()

		r, _ := utils.ToJsValue(c, R)
		return r
	})
}

func JsArgsToGoArgs(info *v8.FunctionCallbackInfo, ctx *v8.Context) []any {
	var args []any
	for _, arg := range info.Args()[1:] {
		v, _ := utils.ToGoValue(ctx, arg)
		args = append(args, v)
	}
	return args
}

func transactionExec(ctx *icode.Ctx, iso *v8.Isolate) *v8.FunctionTemplate {
	return v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) (r *v8.Value) {
		c := info.Context()
		if len(info.Args()) < 1 {
			return utils.JsException(c, "no sql found")
		}
		args := JsArgsToGoArgs(info, c)
		eng := ctx.GetEngine()
		sess := eng.NewSession()
		defer func(sess *xorm.Session) {
			_ = sess.Close()
		}(sess)
		if err := sess.Begin(); err != nil {
			return utils.JsError(c, "error begin transaction")
		}
		if _, err := sess.Exec(info.Args()[0].String(), args); err != nil {
			return utils.JsError(c, "error exec sql")
		}
		err := sess.Commit()
		if err != nil {
			return utils.JsError(c, "error commit transaction")
		}
		r, _ = v8.NewValue(iso, true)
		return
	})
}
func queryInterface(ctx *icode.Ctx, iso *v8.Isolate) *v8.FunctionTemplate {
	return v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) (r *v8.Value) {
		c := info.Context()
		if len(info.Args()) < 1 {
			return utils.JsException(c, "no sql found")
		}
		eng := ctx.GetEngine()
		results, err := eng.QueryInterface(info.Args()[0].String())
		if err != nil {
			return utils.JsException(c, err.Error())
		}
		r, err = utils.ToJsValue(c, results)
		if err != nil {
			return utils.JsError(c, "error convert result to js value")
		}
		return
	})
}
