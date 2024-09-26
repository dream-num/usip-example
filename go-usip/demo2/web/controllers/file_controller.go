// file: controllers/user_controller.go

package controllers

import (
	"go-usip/datamodels"
	"go-usip/services"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"github.com/spf13/viper"
)

type FileController struct {
	// context is auto-binded by Iris on each request,
	// remember that on each incoming request iris creates a new UserController each time,
	// so all fields are request-scoped by-default, only dependency injection is able to set
	// custom fields like the Service which is the same for all requests (static binding)
	// and the Session which depends on the current context (dynamic binding).
	Ctx iris.Context

	Service services.FileService

	// Session, binded using dependency injection from the main.go.
	Session *sessions.Session
}

func (c *FileController) GetList() mvc.Result {
	userId, ok := isLoggedIn(c.Session)
	if !ok {
		return mvc.Response{
			Code: iris.StatusUnauthorized,
			Path: "/user/login",
		}
	}

	files, _ := c.Service.GetByUserId(userId)
	view := mvc.View{
		Name: "file/list.html",
		Data: iris.Map{
			"Files":     files,
			"DocHost":   viper.GetString("univer.docHost"),
			"SheetHost": viper.GetString("univer.sheetHost"),
		},
	}

	return view
}

func (c *FileController) PostNew() mvc.Result {
	userId, ok := isLoggedIn(c.Session)
	if !ok {
		return mvc.Response{
			Code: iris.StatusUnauthorized,
			Path: "/user/login",
		}
	}

	name := c.Ctx.FormValue("name")
	unitType := c.Ctx.FormValue("type")
	file, err := c.Service.Create(services.CreateUnitRequest{
		Name:   name,
		Type:   unitType,
		UserId: userId,
		Cookie: c.Ctx.GetHeader("Cookie"),
	})
	if err != nil {
		return mvc.Response{
			Code: iris.StatusInternalServerError,
			Text: err.Error(),
		}
	}

	path := "/file/list"
	switch datamodels.FileTypeInt(unitType) {
	case datamodels.UnitTypeDoc:
		path = viper.GetString("univer.docHost") + "/?type=1&unit=" + file.UnitId
	case datamodels.UnitTypeSheet:
		path = viper.GetString("univer.sheetHost") + "/?type=2&unit=" + file.UnitId
	}

	return mvc.Response{
		Code: iris.StatusFound,
		Path: path,
	}
}
