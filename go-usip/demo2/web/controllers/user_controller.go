// file: controllers/user_controller.go

package controllers

import (
	"go-usip/datamodels"
	"go-usip/services"
	"image/png"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
)

// UserController is our /user controller.
// UserController is responsible to handle the following requests:
// GET  			/user/register
// POST 			/user/register
// GET 				/user/login
// POST 			/user/login
// All HTTP Methods /user/logout
type UserController struct {
	// context is auto-binded by Iris on each request,
	// remember that on each incoming request iris creates a new UserController each time,
	// so all fields are request-scoped by-default, only dependency injection is able to set
	// custom fields like the Service which is the same for all requests (static binding)
	// and the Session which depends on the current context (dynamic binding).
	Ctx iris.Context

	// Our UserService, it's an interface which
	// is binded from the main application.
	Service services.UserService

	// Session, binded using dependency injection from the main.go.
	Session *sessions.Session
}

func (c *UserController) logout() {
	c.Session.Destroy()
}

// GetRegister handles GET: http://localhost:8080/user/register.
func (c *UserController) GetRegister() mvc.Result {
	return mvc.Response{Path: "/register"}
}

// PostRegister handles POST: http://localhost:8080/user/register.
func (c *UserController) PostRegister() mvc.Result {
	// get nickname, username and password from the form.
	var (
		nickname = c.Ctx.FormValue("nickname")
		username = c.Ctx.FormValue("username")
		password = c.Ctx.FormValue("password")
	)

	// create the new user, the password will be hashed by the service.
	u, err := c.Service.Create(password, datamodels.User{
		Username: username,
		Nickname: nickname,
	})

	// set the user's id to this session even if err != nil,
	// the zero id doesn't matters because .getCurrentUserID() checks for that.
	// If err != nil then it will be shown, see below on mvc.Response.Err: err.
	c.Session.Set(userIDKey, u.UserId)

	return mvc.Response{
		// if not nil then this error will be shown instead.
		Err: err,
		// redirect to files page.
		Path: "/files",
		// When redirecting from POST to GET request you -should- use this HTTP status code,
		// however there're some (complicated) alternatives if you
		// search online or even the HTTP RFC.
		// Status "See Other" RFC 7231, however iris can automatically fix that
		// but it's good to know you can set a custom code;
		// Code: 303,
	}
}

// GetLogin handles GET: http://localhost:8080/user/login.
func (c *UserController) GetLogin() mvc.Result {
	return mvc.Response{Path: "/login"}
}

// PostLogin handles POST: http://localhost:8080/user/register.
func (c *UserController) PostLogin() mvc.Result {
	var (
		username = c.Ctx.FormValue("username")
		password = c.Ctx.FormValue("password")
	)

	u, found := c.Service.GetByUsernameAndPassword(username, password)

	if !found {
		return mvc.Response{
			Path: "/register",
		}
	}

	c.Session.Set(userIDKey, u.UserId)

	return mvc.Response{
		Path: "/files",
	}
}

// GetMe handles GET: http://localhost:8080/user/me.
func (c *UserController) GetMe() mvc.Result {
	return mvc.Response{Path: "/files"}
}

// AnyLogout handles All/Any HTTP Methods for: http://localhost:8080/user/logout.
func (c *UserController) AnyLogout() {
	if _, ok := isLoggedIn(c.Session); ok {
		c.logout()
	}

	c.Ctx.Redirect("/login")
}

func (c *UserController) GetPeople() mvc.Result {
	_, ok := isLoggedIn(c.Session)
	if !ok {
		// if it's not logged in then redirect user to the login page.
		return mvc.Response{Path: "/login"}
	}

	nextId := c.Ctx.URLParamIntDefault("next", 0)
	size := c.Ctx.URLParamIntDefault("size", 10)
	users, latest := c.Service.GetByPage(uint(nextId), uint(size))
	nextId = 0
	if !latest {
		nextId = int(users[len(users)-1].ID)
	}

	c.Ctx.JSON(iris.Map{
		"users": users,
		"next":  nextId,
	})
	return nil
}

func (c *UserController) GetAvatarBy(userId string) mvc.Result {
	// get the avatar by the user ID.
	avatar, found := c.Service.GetAvatarByUserID(userId)
	if !found {
		// if not found then return a simple 404 not found response.
		return mvc.Response{
			Code: iris.StatusNotFound,
		}
	}

	c.Ctx.Header("Content-Type", "image/png")
	c.Ctx.Header("Cache-Control", "public, max-age=604800")
	// write the image to the response.
	png.Encode(c.Ctx.ResponseWriter(), avatar)
	return nil
}
