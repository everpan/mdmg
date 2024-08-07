# v8 模块中几个不支持的特性


1. 函数中`arguments`不支持
2. `iife`之外`let`定义变量会报错：`Identifier 'accept' has already been declared`
   可以采用`var`定义变量，目前未碰到报错