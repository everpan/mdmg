package v8runtime

import (
	"github.com/everpan/mdmg/utils"
	v8 "rogchap.com/v8go"
	"xorm.io/xorm"
)

func ExportXormObject(orm *xorm.Engine, iso *v8.Isolate) *v8.ObjectTemplate {
	ormObj := v8.NewObjectTemplate(iso)
	ormObj.Set("exec", execSql(orm, iso))
	ormObj.Set("tranExec", transactionExec(orm, iso))
	// parent
	obj := v8.NewObjectTemplate(iso)
	obj.Set("sql", ormObj)
	return obj
}

func execSql(orm *xorm.Engine, iso *v8.Isolate) *v8.FunctionTemplate {
	return v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		ctx := info.Context()
		if len(info.Args()) < 0 {
			return utils.JsException(ctx, "no sql found")
		}
		args := JsArgsToGoArgs(info, ctx)
		ret, err := orm.Exec(info.Args()[0].String(), args)
		if err != nil {
			return utils.JsError(ctx, err.Error())
		}
		var R struct {
			LastInsertId int64 `json:"last_insert_id""`
			RowsAffected int64 `json:"rows_affected""`
		}
		R.RowsAffected, _ = ret.LastInsertId()
		R.LastInsertId, _ = ret.RowsAffected()

		r, _ := utils.ToJsValue(ctx, R)
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

func transactionExec(orm *xorm.Engine, iso *v8.Isolate) *v8.FunctionTemplate {
	return v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) (r *v8.Value) {
		ctx := info.Context()
		if len(info.Args()) < 1 {
			return utils.JsException(ctx, "no sql found")
		}
		args := JsArgsToGoArgs(info, ctx)
		sess := orm.NewSession()
		defer sess.Close()
		if err := sess.Begin(); err != nil {
			return utils.JsError(ctx, "error begin transaction")
		}
		if _, err := sess.Exec(info.Args()[0].String(), args); err != nil {
			return utils.JsError(ctx, "error exec sql")
		}
		err := sess.Commit()
		if err != nil {
			return utils.JsError(ctx, "error commit transaction")
		}
		r, _ = v8.NewValue(iso, true)
		return
	})
}
