// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"dksv-v2/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/api/test/", &controllers.MainController{})

	beego.Router("/api/host/info/", &controllers.MainController{}, "get:SysInfo")

	// 容器操作
	beego.Router("/api/container/create/", &controllers.ContainerController{}, "post:Create")
	beego.Router("/api/container/start/", &controllers.ContainerController{}, "post:Start")
	beego.Router("/api/container/inspect/", &controllers.ContainerController{}, "get:Inspect")
	beego.Router("/api/container/stop/", &controllers.ContainerController{}, "post:Stop")
	beego.Router("/api/container/remove/", &controllers.ContainerController{}, "post:Remove")
	beego.Router("/api/container/list/", &controllers.ContainerController{}, "get:List")
	beego.Router("/api/container/logs/", &controllers.ContainerController{}, "get:Logs")
	beego.Router("/api/container/export/", &controllers.ContainerController{}, "post:Export")

	// 镜像操作
	beego.Router("/api/image/list/", &controllers.ImageController{}, "get:List")
	beego.Router("/api/image/pull/", &controllers.ImageController{}, "post:Pull")
	beego.Router("/api/image/remove/", &controllers.ImageController{}, "post:Remove")
	beego.Router("/api/image/push/", &controllers.ImageController{}, "post:Push")
	beego.Router("/api/image/tag/", &controllers.ImageController{}, "post:Tag")

	// 网络操作
	beego.Router("/api/network/create/", &controllers.NetworkController{}, "post:Create")
	beego.Router("/api/network/list/", &controllers.NetworkController{}, "get:List")
	beego.Router("/api/network/remove/", &controllers.NetworkController{}, "post:Remove")

	// 宿主机操作
	beego.Router("/api/host/list/files/", &controllers.HostController{}, "get:ListFiles")
}
