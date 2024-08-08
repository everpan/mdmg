package handler

var ICoderHandler = PathHandler{
	Path:    "/v1/icode/:modVer/:jsFile/*",
	Handler: iCoderHandler,
}
