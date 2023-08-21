package commands

import (
	"encoding/gob"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"lushen/models/mysql"
	"net/url"
	"os"
	"strings"
	"time"
)

func RegisterDataBase() {
	logs.Info("正在初始化数据库配置.")
	dbadapter, _ := web.AppConfig.String("db_adapter")
	orm.DefaultTimeLoc = time.Local
	orm.DefaultRowsLimit = -1

	if strings.EqualFold(dbadapter, "mysql") {
		host, _ := web.AppConfig.String("db_host")
		database, _ := web.AppConfig.String("db_database")
		username, _ := web.AppConfig.String("db_username")
		password, _ := web.AppConfig.String("db_password")

		timezone, _ := web.AppConfig.String("timezone")
		location, err := time.LoadLocation(timezone)
		if err == nil {
			orm.DefaultTimeLoc = location
		} else {
			logs.Error("加载时区配置信息失败,请检查是否存在 ZONEINFO 环境变量->", err)
		}

		port, _ := web.AppConfig.String("db_port")

		dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=%s", username, password, host, port, database, url.QueryEscape(timezone))

		if err := orm.RegisterDataBase("default", "mysql", dataSource); err != nil {
			logs.Error("注册默认数据库失败->", err)
			os.Exit(1)
		}

	} else {
		logs.Error("不支持的数据库类型.")
		os.Exit(1)
	}

	logs.Info("数据库初始化完成.")
}

func RegisterModel() {
	orm.RegisterModel(
		new(mysql.Company),
	)
	gob.Register(mysql.Company{})
}
