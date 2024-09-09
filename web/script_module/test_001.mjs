// import {e} from "./import"
(() => {
    let ctx = __ic.ctx
    let accept = ctx.header().Accept
    return {
        output: {
            code: 0,
            data: {
                sql: "select * from user"
            },
            header: ctx.header("content-type"),
            headers: ctx.header(),
            query: ctx.query('key'),
            queries: ctx.query(),
            param: ctx.param("module"),
            params: ctx.param(),
            accept,
            base: ctx.baseURL(),
            originURL: ctx.originURL()
        }
    }
})()