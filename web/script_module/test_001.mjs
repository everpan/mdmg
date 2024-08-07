// import {e} from "./import"
(() => {
    let accept = icode.header().Accept

    return {
        output: {
            code: 0,
            data: {
                sql: "select * from user"
            },
            header: icode.header("content-type"),
            headers: icode.header(),
            query: icode.query('key'),
            queries: icode.query(),
            param: icode.param("module"),
            params: icode.param(),
            accept,
            base: icode.baseURL(),
            originURL: icode.originURL()
        }
    }
})()