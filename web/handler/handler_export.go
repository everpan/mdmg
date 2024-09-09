package handler

var ICoderHandler = MyHandlerExport{
	Path:    "/v1/icode/:modVer/:jsFile/*",
	Handler: icodeHandler,
}
