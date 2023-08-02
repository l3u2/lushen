package controllers

import (
	"github.com/beego/beego/v2/server/web/context"
)

type UserController struct{}

func (user UserController) AddUser(ctx *context.Context) {
	ctx.WriteString("add user")
}
