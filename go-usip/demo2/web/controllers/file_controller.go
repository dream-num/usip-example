// file: controllers/user_controller.go

package controllers

import (
	"go-usip/datamodels"
	"go-usip/services"
	"io"

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
			"UserId":    userId,
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

func (c *FileController) PostImport() mvc.Result {
	userId, ok := isLoggedIn(c.Session)
	if !ok {
		return mvc.Response{
			Code: iris.StatusUnauthorized,
		}
	}

	unitType := datamodels.FileTypeInt(c.Ctx.FormValue("type"))
	formfile, fileHeader, err := c.Ctx.FormFile("file")
	if err != nil {
		return mvc.Response{
			Code: iris.StatusBadRequest,
			Text: err.Error(),
		}
	}

	file, err := c.Service.Import(services.ImportReq{
		FormFile: formfile,
		UserId:   userId,
		FileName: fileHeader.Filename,
		FileSize: int(fileHeader.Size),
		Type:     unitType,
		Cookie:   c.Ctx.GetHeader("Cookie"),
	})
	if err != nil {
		return mvc.Response{
			Code: iris.StatusInternalServerError,
			Text: err.Error(),
		}
	}

	path := "/file/list"
	switch unitType {
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

func (c *FileController) GetExport() mvc.Result {
	userId, ok := isLoggedIn(c.Session)
	if !ok {
		return mvc.Response{
			Code: iris.StatusUnauthorized,
		}
	}

	fileId, err := c.Ctx.URLParamInt("fileId")
	if err != nil {
		return mvc.Response{
			Code: iris.StatusBadRequest,
			Text: err.Error(),
		}
	}

	result, err := c.Service.Export(services.ExportReq{
		FileId: uint(fileId),
		UserId: userId,
		Cookie: c.Ctx.GetHeader("Cookie"),
	})
	if err != nil {
		return mvc.Response{
			Code: iris.StatusInternalServerError,
			Text: err.Error(),
		}
	}

	c.Ctx.Header("Content-Disposition", "attachment; filename="+result.FileName)
	_, err = io.Copy(c.Ctx.ResponseWriter(), result.Reader)
	if err != nil {
		return mvc.Response{
			Code: iris.StatusInternalServerError,
			Text: err.Error(),
		}
	}
	defer result.Reader.Close()

	return mvc.Response{}
}

func (c *FileController) Delete() mvc.Result {
	userId, ok := isLoggedIn(c.Session)
	if !ok {
		return mvc.Response{
			Code: iris.StatusUnauthorized,
		}
	}

	var req struct {
		FileIds []uint `json:"fileIds"`
	}

	if err := c.Ctx.ReadForm(&req); err != nil {
		return mvc.Response{
			Code: iris.StatusBadRequest,
			Text: err.Error(),
		}
	}

	err := c.Service.BatchDelete(userId, req.FileIds)
	if err != nil {
		return mvc.Response{
			Code: iris.StatusInternalServerError,
			Text: err.Error(),
		}
	}

	return mvc.Response{}
}

func (c *FileController) PostJoin() mvc.Result {
	userId, ok := isLoggedIn(c.Session)
	if !ok {
		return mvc.Response{
			Code: iris.StatusUnauthorized,
		}
	}

	var req struct {
		UserIds []string `json:"userIds"`
		FileId  uint     `json:"fileId"`
		Role    string   `json:"role"`
	}
	if err := c.Ctx.ReadJSON(&req); err != nil {
		return mvc.Response{
			Code: iris.StatusBadRequest,
			Text: err.Error(),
		}
	}

	if !c.Service.CheckPermission(services.CheckPermissionReq{
		FileId: req.FileId,
		UserId: userId,
		Action: services.ActionJoin,
	}) {
		return mvc.Response{
			Code: iris.StatusForbidden,
		}
	}

	err := c.Service.Join(services.JoinReq{
		FileId:  req.FileId,
		Role:    datamodels.Role(req.Role),
		UserIds: req.UserIds,
	})
	if err != nil {
		return mvc.Response{
			Code: iris.StatusInternalServerError,
			Text: err.Error(),
		}
	}

	return mvc.Response{}
}
