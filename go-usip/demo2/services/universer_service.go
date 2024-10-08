package services

import (
	"encoding/json"
	"fmt"
	"go-usip/datamodels"
	"io"
	"log"
	"net/url"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
)

const universerSuccessCode = 1

type UniverserService interface {
	CreateUnit(req CreateUnitRequest) (unitId string, err error)
	UploadFile(req ImportReq) (fileId string, err error)
	Import(req UniverserImportReq) (taskId string, err error)
	PullResult(req UniverserPullReq) (string, error)
	Export(req UniverserExportReq) (taskId string, err error)
	GetFile(req UniverserGetFileReq) (reader io.ReadCloser, err error)
}

func NewUniverseService() UniverserService {
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
		log.Printf("Error while import: %v", err)
		return "", err
	}

	if resp.StatusCode() != 200 {
		log.Printf("Error while import: %v", resp.String())
		return "", fmt.Errorf("Error while import: %v", resp.String())
	}
	var importResp struct {
		Error  UniverserErr `json:"error"`
		TaskId string       `json:"taskID"`
	}
	body := resp.Body()
	if err := json.Unmarshal(body, &importResp); err != nil {
		log.Printf("Error while import: %v", err)
		return "", err
	}

	if importResp.Error.Code != universerSuccessCode {
		log.Printf("Error while import: %v", importResp.Error.Message)
		return "", fmt.Errorf("Error while import: %v", importResp.Error.Message)
	}

	return importResp.TaskId, nil
}

const (
	ExchangeTypeImport = 0
	ExchangeTypeExport = 1
)

type UniverserPullReq struct {
	TaskId       string
	Cookie       string
	ExchangeType int
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
		Export struct {
			FileId string `json:"fileID"`
		} `json:"export"`
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

	log.Printf("Pull result header: %+v", resp.Header())
	log.Printf("Pull result: %+v\n", result)

	switch result.Status {
	case "done":
		if req.ExchangeType == ExchangeTypeImport {
			return result.Import.UnitId, nil
		}
		return result.Export.FileId, nil
	case "pending":
		return "", nil
	default:
		return "", fmt.Errorf("Error while pulling result: %v", result.Status)
	}
}

type UniverserExportReq struct {
	UnitId string
	Type   int
	Cookie string
}

func (s *universeService) Export(req UniverserExportReq) (taskId string, err error) {
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Cookie", req.Cookie).
		SetBody(map[string]interface{}{
			"unitID": req.UnitId,
			"type":   req.Type,
		}).
		Post(fmt.Sprintf("%s/universer-api/exchange/%d/export", viper.GetString("universer.host"), req.Type))
	if err != nil {
		log.Printf("Error while exporting unit: %v", err)
		return
	}
	if resp.StatusCode() != 200 {
		log.Printf("Error while export: %v", resp.String())
		return "", fmt.Errorf("Error while export: %v", resp.String())
	}
	var exportResp struct {
		Error  UniverserErr `json:"error"`
		TaskId string       `json:"taskID"`
	}
	body := resp.Body()
	if err := json.Unmarshal(body, &exportResp); err != nil {
		log.Printf("Error while export: %v", err)
		return "", err
	}

	if exportResp.Error.Code != universerSuccessCode {
		log.Printf("Error while export: %v", exportResp.Error.Message)
		return "", fmt.Errorf("Error while export: %v", exportResp.Error.Message)
	}

	return exportResp.TaskId, nil
}

type UniverserGetFileReq struct {
	FileId string
	Cookie string
}

func (s *universeService) GetFile(req UniverserGetFileReq) (reader io.ReadCloser, err error) {
	client := resty.New()
	resp, err := client.R().
		SetHeader("Cookie", req.Cookie).
		Get(fmt.Sprintf("%s/universer-api/file/%s/sign-url", viper.GetString("universer.host"), req.FileId))
	if err != nil {
		log.Printf("Error while getting file: %v", err)
		return
	}

	if resp.StatusCode() != 200 {
		log.Printf("Error while getting file: %v", resp.String())
		return nil, fmt.Errorf("Error while getting file: %v", resp.String())
	}

	var urlResp struct {
		Error UniverserErr `json:"error"`
		URL   string       `json:"url"`
	}

	body := resp.Body()
	if err := json.Unmarshal(body, &urlResp); err != nil {
		log.Printf("Error while getting file: %v", err)
		return nil, err
	}

	if urlResp.Error.Code != universerSuccessCode {
		log.Printf("Error while getting file: %v", urlResp.Error.Message)
		return nil, fmt.Errorf("Error while getting file: %v", urlResp.Error.Message)
	}

	fileUrl := urlResp.URL

	uri, err := url.Parse(urlResp.URL)
	if err != nil {
		log.Printf("Error while getting file: %v", err)
		return nil, err
	}
	if uri.Host == "" {
		fileUrl = fmt.Sprintf("%s%s", viper.GetString("universer.host"), urlResp.URL)
	}

	resp, err = client.R().
		SetHeader("Cookie", req.Cookie).
		SetDoNotParseResponse(true).
		Get(fileUrl)
	if err != nil {
		log.Printf("Error while getting file: %v", err)
		return
	}

	if resp.StatusCode() != 200 {
		log.Printf("Error while getting file: %v", resp.String())
		return nil, fmt.Errorf("Error while getting file: %v", resp.String())
	}

	reader = resp.RawBody()

	return reader, nil
}
