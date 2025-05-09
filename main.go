package main

import "monitor-server/cmd"

/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2023/10/31 16:29:28
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2023/10/31 16:29:28
 */

// @title monitor-server API
// @version 0.0.1
// @description monitor-server接口文档
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @host localhost:8080
// @BasePath
func main() {
	cmd.Execute()
}
