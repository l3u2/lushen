package admin

import (
	"github.com/beego/beego/v2/core/logs"
	"lushen/models/mysql"
)

type CompanyController struct {
	BaseAdminController
}

func (c *CompanyController) Show() {

	company, err := mysql.NewCompany().Find(1)
	if err != nil {
		logs.Error("获取企业信息失败：", err.Error())
	}

	c.JsonResult(200, "success", company)
}

func (c *CompanyController) Add() {
	m := mysql.NewCompany()
	m.Agzgb = "aaa"
	m.Insert()
}

func (c *CompanyController) Create() {
	m := mysql.NewEmployee()
	m.Create()
}

func (c *CompanyController) Delete() {

}
