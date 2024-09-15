get = () => {
    // 放置外围，会导致变量不允许重复声明错误
    let ctx = __ic.ctx
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