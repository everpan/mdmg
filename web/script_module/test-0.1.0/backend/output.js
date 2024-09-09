get = () => {
    return {
        output: { // 输出结果
            headers: ic.ctx.header(),
            module: ic.ctx.module(),
            version: ic.ctx.version()
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