package services

import (
	"encoding/json"
	"fmt"
	"go-usip/datamodels"
	"log"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
)

const universerSuccessCode = 1

type UniverseService interface {
	CreateUnit(req CreateUnitRequest) (unitId string, err error)
	UploadFile(req ImportReq) (fileId string, err error)
	Import(req UniverserImportReq) (taskId string, err error)
	PullResult(req UniverserPullReq) (string, error)
}

func NewUniverseService() UniverseService {
	return &universeService{}
}

type universeService struct{}

type UniverserErr struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type CreateUnitRequest struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	UserId string `json:"user_id"`

	Cookie string `json:"-"`
}

func (s *universeService) CreateUnit(req CreateUnitRequest) (string, error) {
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Cookie", req.Cookie).
		SetBody(map[string]string{
			"name":    req.Name,
			"creator": req.UserId,
		}).
		Post(fmt.Sprintf("%s/universer-api/snapshot/%d/unit/-/create", viper.GetString("universer.host"), datamodels.FileTypeInt(req.Type)))

	if err != nil {
		log.Printf("Error while creating unit: %v", err)
		return "", err
	}

	if resp.StatusCode() != 200 {
		log.Printf("Error while creating unit: %v", resp.String())
		return "", fmt.Errorf("Error while creating unit: %v", resp.String())
	}
	body := resp.Body()
	var unit struct {
		Error  UniverserErr `json:"error"`
		UnitId string       `json:"unitID"`
	}
	if err := json.Unmarshal(body, &unit); err != nil {
		log.Printf("Error while creating unit: %v", err)
		return "", err
	}

	if unit.Error.Code != universerSuccessCode {
		log.Printf("Error while creating unit: %v", unit.Error.Message)
		return "", fmt.Errorf("Error while creating unit: %v", unit.Error.Message)
	}

	return unit.UnitId, nil
}

type UploadFileResp struct {
	FileId string `json:"FileId"`
}

func (s *universeService) UploadFile(req ImportReq) (fileId string, err error) {
	client := resty.New()
	resp, err := client.R().
		SetHeader("Cookie", req.Cookie).
		SetFileReader("file", req.FileName, req.FormFile).
		Post(fmt.Sprintf("%s/universer-api/stream/file/upload?size=%d", viper.GetString("universer.host"), req.FileSize))
	if err != nil {
		log.Printf("Error while uploading file: %v", err)
		return
	}

	if resp.StatusCode() != 200 {
		log.Printf("Error while uploading file: %v", resp.String())
		return "", fmt.Errorf("Error while uploading file: %v", resp.String())
	}

	body := resp.Body()
	var fileResp struct {
		FileId string `json:"FileId"`
	}
	if err = json.Unmarshal(body, &fileResp); err != nil {
		log.Printf("Error while uploading file: %v", err)
		return
	}

	return fileResp.FileId, nil
}

type UniverserImportReq struct {
	FileId     string
	Type       int
	OutputType int

	Cookie string `json:"-"`
}

func (s *universeService) Import(req UniverserImportReq) (taskId string, err error) {
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Cookie", req.Cookie).
		SetBody(map[string]interface{}{
			"fileID":     req.FileId,
			"outputType": req.OutputType,
		}).
		Post(fmt.Sprintf("%s/universer-api/exchange/%d/import", viper.GetString("universer.host"), req.Type))
	if err != nil {
		log.Printf("Error while uploading file: %v", err)
		return "", err
	}

	if resp.StatusCode() != 200 {
		log.Printf("Error while uploading file: %v", resp.String())
		return "", fmt.Errorf("Error while uploading file: %v", resp.String())
	}
	var importResp struct {
		Error  UniverserErr `json:"error"`
		TaskId string       `json:"taskID"`
	}
	body := resp.Body()
	if err := json.Unmarshal(body, &importResp); err != nil {
		log.Printf("Error while uploading file: %v", err)
		return "", err
	}

	if importResp.Error.Code != universerSuccessCode {
		log.Printf("Error while uploading file: %v", importResp.Error.Message)
		return "", fmt.Errorf("Error while uploading file: %v", importResp.Error.Message)
	}

	return importResp.TaskId, nil
}

type UniverserPullReq struct {
	TaskId string
	Cookie string
}

func (s *universeService) PullResult(req UniverserPullReq) (string, error) {
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Cookie", req.Cookie).
		Get(fmt.Sprintf("%s/universer-api/exchange/task/%s", viper.GetString("universer.host"), req.TaskId))
	if err != nil {
		log.Printf("Error while pulling result: %v", err)
		return "", err
	}

	if resp.StatusCode() != 200 {
		log.Printf("Error while pulling result: %v", resp.String())
		return "", fmt.Errorf("Error while pulling result: %v", resp.String())
	}

	var result struct {
		Error  UniverserErr `json:"error"`
		Status string       `json:"status"`
		Import struct {
			UnitId string `json:"unitID"`
		} `json:"import"`
	}
	body := resp.Body()
	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("Error while pulling result: %v", err)
		return "", err
	}

	if result.Error.Code != universerSuccessCode {
		log.Printf("Error while pulling result: %v", result.Error.Message)
		return "", fmt.Errorf("Error while pulling result: %v", result.Error.Message)
	}

	switch result.Status {
	case "done":
		return result.Import.UnitId, nil
	case "pending":
		return "", nil
	default:
		return "", fmt.Errorf("Error while pulling result: %v", result.Status)
	}
}
