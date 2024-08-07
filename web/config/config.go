package config

type Config struct {
	JSModuleRootPath string
}

var DefaultConfig = Config{
	JSModuleRootPath: "./web/script_module",
}
