get = () => {
    return {
        output: { // 输出结果
            headers: icode.header(),
            module: icode.module(),
            version: icode.version()
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