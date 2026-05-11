package controllers

import (
	"go-usip/datamodels"
	"go-usip/services"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
)

type AuthAPIController struct {
	Ctx iris.Context

	Service services.UserService
	Session *sessions.Session
}

type apiErrorResp struct {
	Error string `json:"error"`
}

func writeAPIError(ctx iris.Context, code int, message string) mvc.Result {
	ctx.StatusCode(code)
	ctx.JSON(apiErrorResp{Error: message})
	return nil
}

type authUserResp struct {
	UserId   string `json:"userId"`
	Nickname string `json:"nickname"`
	Username string `json:"username"`
}

type authSuccessResp struct {
	User authUserResp `json:"user"`
}

func buildAuthUserResp(user datamodels.User) authUserResp {
	return authUserResp{
		UserId:   user.UserId,
		Nickname: user.Nickname,
		Username: user.Username,
	}
}

func (c *AuthAPIController) GetMe() mvc.Result {
	userID, ok := isLoggedIn(c.Session)
	if !ok {
		return writeAPIError(c.Ctx, iris.StatusUnauthorized, "unauthorized")
	}

	user, found := c.Service.GetByID(userID)
	if !found {
		c.Session.Destroy()
		return writeAPIError(c.Ctx, iris.StatusUnauthorized, "unauthorized")
	}

	c.Ctx.JSON(authSuccessResp{User: buildAuthUserResp(user)})
	return nil
}

type authLoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c *AuthAPIController) PostLogin() mvc.Result {
	var req authLoginReq
	if err := c.Ctx.ReadJSON(&req); err != nil {
		return writeAPIError(c.Ctx, iris.StatusBadRequest, "invalid request body")
	}

	user, found := c.Service.GetByUsernameAndPassword(req.Username, req.Password)
	if !found {
		return writeAPIError(c.Ctx, iris.StatusUnauthorized, "invalid username or password")
	}

	c.Session.Set(userIDKey, user.UserId)
	c.Ctx.JSON(authSuccessResp{User: buildAuthUserResp(user)})
	return nil
}

type authRegisterReq struct {
	Nickname string `json:"nickname"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c *AuthAPIController) PostRegister() mvc.Result {
	var req authRegisterReq
	if err := c.Ctx.ReadJSON(&req); err != nil {
		return writeAPIError(c.Ctx, iris.StatusBadRequest, "invalid request body")
	}

	user, err := c.Service.Create(req.Password, datamodels.User{
		Nickname: req.Nickname,
		Username: req.Username,
	})
	if err != nil {
		return writeAPIError(c.Ctx, iris.StatusBadRequest, err.Error())
	}

	c.Session.Set(userIDKey, user.UserId)
	c.Ctx.StatusCode(iris.StatusCreated)
	c.Ctx.JSON(authSuccessResp{User: buildAuthUserResp(user)})
	return nil
}

func (c *AuthAPIController) PostLogout() mvc.Result {
	c.Session.Destroy()
	c.Ctx.StatusCode(iris.StatusOK)
	c.Ctx.JSON(iris.Map{"ok": true})
	return nil
}
