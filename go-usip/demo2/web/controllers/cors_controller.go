package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"github.com/spf13/viper"
)

type CorsController struct {
	// context is auto-binded by Iris on each request,
	// remember that on each incoming request iris creates a new UserController each time,
	// so all fields are request-scoped by-default, only dependency injection is able to set
	// custom fields like the Service which is the same for all requests (static binding)
	// and the Session which depends on the current context (dynamic binding).
	Ctx iris.Context

	// Session, binded using dependency injection from the main.go.
	Session *sessions.Session
}

func (c *CorsController) GetPage() mvc.Result {
	c.Ctx.JSON(iris.Map{
		"host": viper.GetString("universer.host"),
	})
	return nil
}
