package config

type Config struct {
	JSModuleRootPath string
	// module-version/${JSModuleBeckEndDir}/script.js
	JSModuleBeckEndDir string
}

var DefaultConfig = Config{
	JSModuleRootPath:   "web/script_module",
	JSModuleBeckEndDir: "backend",
}
