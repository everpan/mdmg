get = () => {
    return {
        output: {
            // tenantList: __ic.db.query("select * from ic_tenant_info"),
            tenant_info: __ic.ctx.tenant()
        }
    }
}
(() => {
    return {
        get
    }
})()