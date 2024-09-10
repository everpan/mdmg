get = () => {
    return {
        output: {
            tenantList: __ic.db.query("select * from ic_tenant_info")
        }
    }
}
(() => {
    return {
        get
    }
})()