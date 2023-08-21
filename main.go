package main

import (
	"flag"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"lushen/commands"
	_ "lushen/routers"
	"time"
)

func main() {
	Port := flag.String("port", "8518", "端口号")
	flag.Parse()

	logs.SetLogger(logs.AdapterConsole)
	dateStr := time.Now().Format("20060102")
	config := `{"filename":"lushendata/logs/atsjkhelper_%s.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10,"color":true}`
	config = fmt.Sprintf(config, dateStr)
	err := logs.SetLogger(logs.AdapterFile, config)
	if err != nil {
		panic(err)
	}

	commands.RegisterDataBase()
	commands.RegisterModel()

	beego.Run(":" + *Port)
}
