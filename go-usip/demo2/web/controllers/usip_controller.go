package controllers

import (
	"fmt"
	"go-usip/services"
	"log"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"github.com/spf13/viper"
)

type UsipController struct {
	// context is auto-binded by Iris on each request,
	// remember that on each incoming request iris creates a new UserController each time,
	// so all fields are request-scoped by-default, only dependency injection is able to set
	// custom fields like the Service which is the same for all requests (static binding)
	// and the Session which depends on the current context (dynamic binding).
	Ctx iris.Context

	UserService services.UserService
	FileService services.FileService

	// Session, binded using dependency injection from the main.go.
	Session *sessions.Session
}

type UsipUser struct {
	UserId string `json:"userID,omitempty"`
	Name   string `json:"name,omitempty"`
	Avatar string `json:"avatar,omitempty"`
}

type UsipCredentialResp struct {
	User UsipUser `json:"user,omitempty"`
}

func (c *UsipController) GetCredential() mvc.Result {
	log.Printf("GetCredential: %+v", c.Ctx.Request().Header)
	userId, ok := isLoggedIn(c.Session)
	if !ok {
		return mvc.Response{
			Code: iris.StatusUnauthorized,
		}
	}

	user, ok := c.UserService.GetByID(userId)
	if !ok {
		return mvc.Response{
			Code: iris.StatusUnauthorized,
		}
	}

	c.Ctx.JSON(UsipCredentialResp{
		User: UsipUser{
			UserId: user.UserId,
			Name:   user.Nickname,
			Avatar: fmt.Sprintf("%s/user/avatar/%s", viper.GetString("host"), user.UserId),
		},
	})

	return nil
}

type UsipUserinfoReq struct {
	UserIds []string `json:"userIDs,omitempty"`
}

type UsipUserinfoResp struct {
	Users []UsipUser `json:"users,omitempty"`
}

func (c *UsipController) PostUserinfo() mvc.Result {
	var req UsipUserinfoReq
	if err := c.Ctx.ReadJSON(&req); err != nil {
		return mvc.Response{
			Code: iris.StatusBadRequest,
		}
	}

	tmp, found := c.UserService.GetInIDs(req.UserIds)
	if !found {
		c.Ctx.JSON(UsipUserinfoResp{Users: []UsipUser{}})
		return nil
	}

	users := make([]UsipUser, 0, len(req.UserIds))
	for _, u := range tmp {
		users = append(users, UsipUser{
			UserId: u.UserId,
			Name:   u.Nickname,
			Avatar: fmt.Sprintf("%s/user/avatar/%s", viper.GetString("host"), u.UserId),
		})
	}

	c.Ctx.JSON(UsipUserinfoResp{Users: users})
	return nil
}

type UsipGetRoleResp struct {
	UserId string `json:"userID,omitempty"`
	Role   string `json:"role,omitempty"`
}

func (c *UsipController) GetRole() mvc.Result {
	userId := c.Ctx.FormValue("userID")
	unitId := c.Ctx.FormValue("unitID")

	collaborators, found := c.FileService.GetCollaboratorsByUnitId(unitId)
	if !found {
		return mvc.Response{
			Code: iris.StatusNotFound,
		}
	}

	for _, v := range collaborators {
		if v.UserId == userId {
			c.Ctx.JSON(UsipGetRoleResp{
				UserId: userId,
				Role:   string(v.Role),
				// Avatar: v.Avatar,
			})
			return nil
		}
	}

	return mvc.Response{
		Code: iris.StatusNotFound,
	}
}

type UsipSubject struct {
	ID     string `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Avatar string `json:"avatar,omitempty"`
	// Type   string `json:"type,omitempty"`
}
type UsipCollaborator struct {
	Subject UsipSubject `json:"subject,omitempty"`
	Role    string      `json:"role,omitempty"`
}
type UsipCollaboratorsResp struct {
	Collaborators []struct {
		UnitId   string             `json:"unitID,omitempty"`
		Subjects []UsipCollaborator `json:"subjects,omitempty"`
	} `json:"collaborators"`
}

type UsipCollaboratorsReq struct {
	UnitIds []string `json:"unitIDs,omitempty"`
}

func (c *UsipController) PostCollaborators() mvc.Result {
	var req UsipCollaboratorsReq
	if err := c.Ctx.ReadJSON(&req); err != nil {
		return mvc.Response{
			Code: iris.StatusBadRequest,
		}
	}

	resp := UsipCollaboratorsResp{}
	for _, unitId := range req.UnitIds {
		collaborators, found := c.FileService.GetCollaboratorsByUnitId(unitId)
		if !found {
			continue
		}

		subjects := make([]UsipCollaborator, 0, len(collaborators))
		for _, v := range collaborators {
			user, found := c.UserService.GetByID(v.UserId)
			if !found {
				continue
			}
			subjects = append(subjects, UsipCollaborator{
				Subject: UsipSubject{
					ID:     user.UserId,
					Name:   user.Nickname,
					Avatar: fmt.Sprintf("%s/user/avatar/%s", viper.GetString("host"), user.UserId),
				},
				Role: string(v.Role),
			})
		}
		resp.Collaborators = append(resp.Collaborators, struct {
			UnitId   string             `json:"unitID,omitempty"`
			Subjects []UsipCollaborator `json:"subjects,omitempty"`
		}{
			UnitId:   unitId,
			Subjects: subjects,
		})
	}

	c.Ctx.JSON(resp)
	return nil
}
