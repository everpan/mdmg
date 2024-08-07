/// iife

var accept = icode.header().Accept // 这里可以使用var定义全局变量，不能使用let定义
console.log(accept)
get = () => {
    return {
        output: {
            code: 0,
            data: {
                sql: "select * from user"
            },
            method: "get",
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
}

post = () => {
    return {
        method: "post",
        output: {
            accept,
        }
    }
}

del = () => {
    return {
        method: "del",
        output: {
            accept,
            headers: icode.header(),
            data: "some data",
            "content-type": icode.header("content-type"),
        }
    }
}
// iife
(() => {
        return {
            get,
            post,
            del,
        }
    }
)();