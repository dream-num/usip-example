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
