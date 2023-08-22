package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"lushen/controllers"
	"lushen/controllers/admin"
)

func init() {
	// 注解路由 官方建议：我们不再推荐使用这种方式，因为可读性和可维护性都不太好。特别是重构进行方法重命名的时候，容易出错。
	beego.Router("/", &controllers.MainController{})
	beego.Router("/company/show", &admin.CompanyController{}, "get:Show")
	beego.Router("/company/add", &admin.CompanyController{}, "get:Add")
	beego.Router("/company/create", &admin.CompanyController{}, "get:Create")
	beego.Router("/company/delete", &admin.CompanyController{}, "get:Delete")

	// 优先使用函数式风格的路由注册
	var user controllers.UserController
	beego.Get("/user/add", user.AddUser)
}
