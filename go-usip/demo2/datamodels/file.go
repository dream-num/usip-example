package datamodels

import (
	"gorm.io/gorm"
)

const (
	UnitTypeDoc   = 1
	UnitTypeSheet = 2
)

type File struct {
	gorm.Model
	Name     string `json:"name"`
	UnitId   string `json:"unit_id"`
	UnitType int    `json:"unit_type"`
}

func FileTypeStr(unitType int) string {
	switch unitType {
	case UnitTypeDoc:
		return "doc"
	case UnitTypeSheet:
		return "sheet"
	default:
		return "unknown"
	}
}

func FileTypeInt(unitType string) int {
	switch unitType {
	case "doc":
		return UnitTypeDoc
	case "sheet":
		return UnitTypeSheet
	default:
		return 0
	}
}
