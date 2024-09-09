package v8runtime

import (
	"github.com/everpan/mdmg/utils"
	v8 "rogchap.com/v8go"
	"xorm.io/xorm"
)

func ExportXormObject(engine *xorm.Engine, iso *v8.Isolate) *v8.ObjectTemplate {
	obj := v8.NewObjectTemplate(iso)
	_ = obj.Set("exec", execSql(engine, iso))
	_ = obj.Set("transaction_exec", transactionExec(engine, iso))
	_ = obj.Set("tExec", transactionExec(engine, iso))
	_ = obj.Set("query", queryInterface(engine, iso))

	return obj
}

func execSql(engine *xorm.Engine, iso *v8.Isolate) *v8.FunctionTemplate {
	return v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		c := info.Context()
		if len(info.Args()) < 0 {
			return utils.JsException(c, "no sql found")
		}
		args, _ := utils.ToGoValues(c, info.Args())
		ret, err := engine.Exec(args...)
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

func transactionExec(engine *xorm.Engine, iso *v8.Isolate) *v8.FunctionTemplate {
	return v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) (r *v8.Value) {
		c := info.Context()
		if len(info.Args()) < 1 {
			return utils.JsException(c, "no sql found")
		}
		sess := engine.NewSession()
		defer func(sess *xorm.Session) {
			_ = sess.Close()
		}(sess)
		if err := sess.Begin(); err != nil {
			return utils.JsError(c, "error begin transaction")
		}
		args, _ := utils.ToGoValues(c, info.Args())
		if _, err := engine.Exec(args...); err != nil {
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
func queryInterface(engine *xorm.Engine, iso *v8.Isolate) *v8.FunctionTemplate {
	return v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) (r *v8.Value) {
		c := info.Context()
		if len(info.Args()) < 1 {
			return utils.JsException(c, "no sql found")
		}
		args, _ := utils.ToGoValues(c, info.Args())
		results, err := engine.QueryInterface(args...)
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
