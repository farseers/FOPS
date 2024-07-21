package main

import (
	"github.com/farseer-go/fs"
	"github.com/farseer-go/webapi"
)

func main() {
	fs.Initialize[StartupModule]("fops")
	webapi.RegisterRoutes(route...)
	webapi.UseCors()        // 使用CORS中间件
	webapi.UseApiResponse() // 让所有的返回值，包含在core.ApiResponse中
	//webapi.PrintRoute()     // 打印所有路由信息到控制台
	webapi.UseApiDoc()      // 开启api doc文档
	webapi.UseValidate()    // 使用DTO验证
	webapi.UseStaticFiles() // 使用静态文件 在根目录./wwwroot中的文件
	webapi.UseHealthCheck() // 开启健康检查
	webapi.UsePprof()       // 开启pprof
	webapi.Run()            // 运行web服务，端口配置在：farseer.yaml Webapi.Url 配置节点
}
