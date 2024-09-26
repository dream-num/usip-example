package repositories

import (
	"go-usip/datamodels"
	"log"

	"gorm.io/gorm"
)

type FileCollaboratorRepository interface {
	GetByUserId(userId string) ([]datamodels.FileCollaborator, bool)
	GetByFileId(fileId uint) ([]datamodels.FileCollaborator, bool)

	Create(fileCollaborator datamodels.FileCollaborator) (datamodels.FileCollaborator, error)
}

func NewFileCollaboratorRepository(db *gorm.DB) FileCollaboratorRepository {
	if err := db.AutoMigrate(&datamodels.FileCollaborator{}); err != nil {
		log.Fatalf("AutoMigrate error: %v", err)
	}

	return &fileCollaboratorRepository{db: db}
}

type fileCollaboratorRepository struct {
	db *gorm.DB
}

func (r *fileCollaboratorRepository) GetByUserId(userId string) ([]datamodels.FileCollaborator, bool) {
	var fileCollaborators []datamodels.FileCollaborator
	if err := r.db.Where("user_id = ?", userId).Find(&fileCollaborators).Error; err != nil {
		log.Printf("Error while getting collaborator by user_id: %v", err)
		return fileCollaborators, false
	}
	return fileCollaborators, true
}

func (r *fileCollaboratorRepository) GetByFileId(fileId uint) ([]datamodels.FileCollaborator, bool) {
	var fileCollaborators []datamodels.FileCollaborator
	if err := r.db.Where("file_id = ?", fileId).Find(&fileCollaborators).Error; err != nil {
		log.Printf("Error while getting collaborator by file_id: %v", err)
		return fileCollaborators, false
	}
	return fileCollaborators, true
}

func (r *fileCollaboratorRepository) Create(fileCollaborator datamodels.FileCollaborator) (datamodels.FileCollaborator, error) {
	return fileCollaborator, r.db.Create(&fileCollaborator).Error
}
