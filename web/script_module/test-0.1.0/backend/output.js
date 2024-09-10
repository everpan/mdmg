let ctx = __ic.ctx
get = () => {
    return {
        output: { // 输出结果
            headers: ctx.header(),
            module: ctx.module(),
            version: ctx.version()
        }
    }
}

del = () => {
    return {
        output: {
            method: "delete is keyword in js, alias is del"
        }
    }
}

patch = () => {
    return {
        msg: "not found output in response"
    }
}

(() => {
    return {
        get, del, patch
    }
})()