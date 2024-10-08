package repositories

import (
	"go-usip/datamodels"
	"log"

	"gorm.io/gorm"
)

type FileRepository interface {
	Get(id uint) (file datamodels.File, found bool)
	GetByUnitId(unitId string) (datamodels.File, bool)
	BatchGet(ids []uint) (files []datamodels.File, found bool)

	Create(file datamodels.File) (datamodels.File, error)

	BatchDelete(ids []uint) error
}

func NewFileRepository(db *gorm.DB) FileRepository {
	if err := db.AutoMigrate(&datamodels.File{}); err != nil {
		log.Fatalf("AutoMigrate error: %v", err)
	}

	return &fileRepository{db: db}
}

type fileRepository struct {
	db *gorm.DB
}

func (r *fileRepository) Get(id uint) (file datamodels.File, found bool) {
	file = datamodels.File{}
	if err := r.db.Where("id = ?", id).First(&file).Error; err != nil {
		log.Printf("Error while getting file by id: %v", err)
		return file, false
	}
	return file, true
}

func (r *fileRepository) GetByUnitId(unitId string) (datamodels.File, bool) {
	file := datamodels.File{}
	if err := r.db.Where("unit_id = ?", unitId).First(&file).Error; err != nil {
		log.Printf("Error while getting file by id: %v", err)
		return file, false
	}
	return file, true
}

func (r *fileRepository) BatchGet(ids []uint) (files []datamodels.File, found bool) {
	if err := r.db.Where("id IN ?", ids).Find(&files).Error; err != nil {
		log.Printf("Error while getting files by ids: %v", err)
		return files, false
	}
	return files, true
}

func (r *fileRepository) Create(file datamodels.File) (datamodels.File, error) {
	return file, r.db.Create(&file).Error
}

func (r *fileRepository) BatchDelete(ids []uint) error {
	return r.db.Where("id IN ?", ids).Delete(&datamodels.File{}).Error
}
