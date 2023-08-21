package mysql

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type Company struct {
	Id        int       `json:"id" orm:"pk;auto;unique;column(id);description(主键)"`
	Bk        string    `json:"bk" orm:"column(bk);size(16);description(板块)"`
	Agdm      string    `json:"agdm" orm:"column(agdm);size(16);description(股票代码)"`
	Agjc      string    `json:"agjc" orm:"column(agjc);size(32);description(A股简称)"`
	Agssrq    string    `json:"agssrq" orm:"column(agssrq);size(32);description(上市日期)"`
	Agzgb     string    `json:"agzgb" orm:"column(agzgb);size(32);description(总股本)"`
	Agltgb    string    `json:"agltgb" orm:"column(agltgb);size(32);description(流通股本)"`
	Sshymc    string    `json:"sshymc" orm:"column(sshymc);size(32);description(所属行业)"`
	CreatedAt time.Time `json:"created_at" orm:"column(created_at);type(datetime);auto_now_add;description(创建时间)"`
	UpdatedAt time.Time `json:"updated_at" orm:"column(updated_at);type(datetime);auto_now;description(最后修改时间)"`
}

func (c *Company) TableName() string {
	return "company"
}

func (c *Company) TableEngine() string {
	return "INNODB"
}

func NewCompany() *Company {
	return &Company{}
}

func (c *Company) Find(id int) (*Company, error) {
	o := orm.NewOrm()

	err := o.QueryTable(c.TableName()).Filter("id", id).One(c)
	if err != nil {
		logs.Error("查询文章时失败 -> ", err)
		return nil, err
	}

	return c, nil
}
