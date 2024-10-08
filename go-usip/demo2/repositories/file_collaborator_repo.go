package repositories

import (
	"go-usip/datamodels"
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type FileCollaboratorRepository interface {
	Get(fileId uint, userId string) (datamodels.FileCollaborator, bool)
	GetByUserId(userId string) ([]datamodels.FileCollaborator, bool)
	GetByFileId(fileId uint) ([]datamodels.FileCollaborator, bool)

	Create(fileCollaborator datamodels.FileCollaborator) (datamodels.FileCollaborator, error)
	InsertOrUpdate(fileCollaborators []datamodels.FileCollaborator) error

	BatchDelete(userId string, fileIds []uint) error
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

func (r *fileCollaboratorRepository) Get(fileId uint, userId string) (datamodels.FileCollaborator, bool) {
	var fileCollaborator datamodels.FileCollaborator
	if err := r.db.Where("file_id = ? AND user_id = ?", fileId, userId).First(&fileCollaborator).Error; err != nil {
		log.Printf("Error while getting collaborator: %v", err)
		return fileCollaborator, false
	}
	return fileCollaborator, true
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

func (r *fileCollaboratorRepository) InsertOrUpdate(fileCollaborators []datamodels.FileCollaborator) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "file_id"}, {Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"role"}),
	}).Create(&fileCollaborators).Error
}

func (r *fileCollaboratorRepository) BatchDelete(userId string, fileIds []uint) error {
	return r.db.Where("user_id = ? AND file_id IN ?", userId, fileIds).Delete(&datamodels.FileCollaborator{}).Error
}
