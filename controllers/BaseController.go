package controllers

import (
	"encoding/json"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"io"
	"time"
)

type BaseController struct {
	web.Controller
}

func (c *BaseController) JsonResult(code int, message string, data ...interface{}) {
	jsonData := make(map[string]interface{}, 4)

	jsonData["code"] = code
	jsonData["message"] = message
	jsonData["serverTime"] = time.Now().Format("2006-01-02 15:04:05")

	if len(data) > 0 && data[0] != nil {
		jsonData["data"] = data[0]
	}

	returnJSON, err := json.Marshal(jsonData)
	if err != nil {
		logs.Error(err)
	}

	c.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	c.Ctx.ResponseWriter.Header().Set("Cache-Control", "no-cache, no-store")
	_, err = io.WriteString(c.Ctx.ResponseWriter, string(returnJSON))
	if err != nil {
		logs.Error(err)
	}

	c.StopRun()
}
