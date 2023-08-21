package admin

import (
	"fmt"
	"github.com/beego/beego/v2/adapter/logs"
	"lushen/models/mysql"
)

type CompanyController struct {
	BaseAdminController
}

func (c *CompanyController) Show() {
	company, err := mysql.NewCompany().Find(1)
	fmt.Println(company)

	if err != nil {
		logs.Error("获取企业信息失败：", err.Error())
	}

	c.JsonResult(200, "success", company)
}
