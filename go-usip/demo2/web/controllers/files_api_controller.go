package controllers

import (
	"go-usip/datamodels"
	"go-usip/services"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"github.com/spf13/viper"
)

type FilesAPIController struct {
	Ctx iris.Context

	Service services.FileService
	Session *sessions.Session
}

type fileItemResp struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	UnitId    string `json:"unitId"`
	UnitType  int    `json:"unitType"`
	Role      string `json:"role"`
	UpdatedAt string `json:"updatedAt"`
	OpenURL   string `json:"openUrl"`
	ExportURL string `json:"exportUrl"`
}

type filesListResp struct {
	UserId string         `json:"userId"`
	Files  []fileItemResp `json:"files"`
}

func (c *FilesAPIController) Get() mvc.Result {
	userID, ok := isLoggedIn(c.Session)
	if !ok {
		return writeAPIError(c.Ctx, iris.StatusUnauthorized, "unauthorized")
	}

	files, _ := c.Service.GetByUserId(userID)
	host := c.Ctx.Host()
	sheetHost := getUnitHost(viper.GetString("univer.sheetHost"), host)

	resp := filesListResp{
		UserId: userID,
		Files:  make([]fileItemResp, 0, len(files)),
	}

	for _, file := range files {
		item := fileItemResp{
			ID:        file.ID,
			Name:      file.Name,
			UnitId:    file.UnitId,
			UnitType:  file.UnitType,
			Role:      string(datamodels.RoleReader),
			UpdatedAt: file.UpdatedAt.Format("2006-01-02 15:04:05"),
			ExportURL: "/file/export?fileId=" + strconv.Itoa(int(file.ID)),
		}

		collaborators, found := c.Service.GetCollaborators(file.ID)
		if found {
			for _, collaborator := range collaborators {
				if collaborator.UserId == userID && collaborator.Role != "" {
					item.Role = string(collaborator.Role)
					break
				}
			}
		}

		if file.UnitType == datamodels.UnitTypeSheet {
			item.OpenURL = sheetHost + "/?unit=" + file.UnitId + "&type=2"
		}

		resp.Files = append(resp.Files, item)
	}

	c.Ctx.JSON(resp)
	return nil
}
