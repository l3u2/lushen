package mysql

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/client/orm/filter/bean"
	"github.com/beego/beego/v2/core/logs"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type Company struct {
	// 注释中禁止包含引号
	Id     uint      `json:"id" orm:"auto;column(id);description(主键)"`
	Bk     string    `json:"bk" orm:"column(bk);size(16);default(主板);description(板块)"`
	Agdm   string    `json:"agdm" orm:"unique;column(agdm);size(16);description(股票代码)"`
	Agjc   string    `json:"agjc" orm:"unique;column(agjc);size(32);description(A股简称)"`
	Agssrq time.Time `json:"agssrq" orm:"index;column(agssrq);type(date);description(上市日期)"`
	Agzgb  string    `json:"agzgb" orm:"index;column(agzgb);size(32);description(总股本)"`
	Agltgb string    `json:"agltgb" orm:"index;column(agltgb);size(32);description(流通股本)"`
	Sshymc string    `json:"sshymc" orm:"index;column(sshymc);size(32);description(所属行业)"`
	// 对于批量的 update 自动更新时间是不生效的
	CreatedAt time.Time `json:"created_at" orm:"column(created_at);type(datetime);auto_now_add;description(创建时间)"`
	UpdatedAt time.Time `json:"updated_at" orm:"column(updated_at);type(datetime);auto_now;description(最后修改时间)"`
}

// 自定义表名
func (c *Company) TableName() string {
	return "company"
}

// 自定义引擎
func (c *Company) TableEngine() string {
	return "INNODB"
}

// 多字段索引
func (c *Company) TableIndex() [][]string {
	return [][]string{
		{"Sshymc", "Agltgb"},
	}
}

// 实例化模型
func NewCompany() *Company {
	return &Company{}
}

// 创建表结构
func (c *Company) Create() bool {
	err := orm.RunSyncdb("default", false, true)
	if err != nil {
		logs.Error("Create Table Error:", err.Error())
		return false
	}

	logs.Info("Create Table Success")
	return true
}

// 自动填充默认值的插入
func (c *Company) Insert() (int64, error) {
	builder := bean.NewDefaultValueFilterChainBuilder(nil, true, true)
	orm.AddGlobalFilterChain(builder.FilterChain)
	o := orm.NewOrm()
	id, err := o.Insert(c)
	if err != nil {
		logs.Error("添加企业失败 -->", err.Error())
		return 0, err
	}

	return id, nil
}

// 更新企业
func (c *Company) Update(cols ...string) (int64, error) {
	o := orm.NewOrm()
	affectedRows, err := o.Update(c, cols...)

	if err != nil {
		logs.Error("更新企业失败 -->", err.Error())
		return 0, err
	}

	return affectedRows, nil
}

// 插入或更新
func (c *Company) InsertOrUpdate() (company *Company, err error) {
	o := orm.NewOrm()

	if c.Id > 0 {
		_, err = o.Update(c)
	} else {
		_, err = o.Insert(c)
	}

	company = c
	return
}

// 根据主键删除企业
func (c *Company) Delete(id uint) (int64, error) {
	o := orm.NewOrm()
	affectedRows, err := o.QueryTable(c.TableName()).Filter("id", id).Delete()

	if err != nil {
		logs.Error("删除id是", id, "的企业失败 -->", err.Error())
		return 0, err
	}

	return affectedRows, nil
}

// 主键查询
func (c *Company) Find(id int, cols ...string) (*Company, error) {
	o := orm.NewOrm()
	err := o.QueryTable(c.TableName()).Filter("id", id).One(c, cols...)
	if err != nil {
		logs.Error("查询企业失败 -> ", err)
		return c, err
	}

	return c, nil
}

// 分页查询
func (c *Company) FindToPager(pageIndex, pageSize int) (companies []*Company, totalCount int, err error) {
	o := orm.NewOrm()

	count, err := o.QueryTable(NewCompany().TableName()).Count()
	if err != nil {
		return
	}
	totalCount = int(count)

	sql := `SELECT e.*,c.agdm
			FROM employee AS e
			INNER JOIN company AS c ON e.company_id = c.id
			ORDER BY c.updated_at DESC limit ? offset ?`

	offset := (pageIndex - 1) * pageSize

	_, err = o.Raw(sql, pageSize, offset).QueryRows(&companies)

	return
}
