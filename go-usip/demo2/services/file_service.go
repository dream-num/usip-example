package services

import (
	"go-usip/datamodels"
	"go-usip/repositories"
	"log"
)

type FileService interface {
	GetByUserId(userId string) ([]datamodels.File, bool)
	GetByFileId(fileId uint) (datamodels.File, bool)
	GetCollaborators(fileId uint) ([]datamodels.FileCollaborator, bool)
	GetCollaboratorsByUnitId(unitId string) ([]datamodels.FileCollaborator, bool)

	Create(req CreateUnitRequest) (datamodels.File, error)
}

type fileService struct {
	repo      repositories.FileRepository
	collaRepo repositories.FileCollaboratorRepository

	uSvc UniverseService
}

func NewFileService(repo repositories.FileRepository, collaRepo repositories.FileCollaboratorRepository, uSvc UniverseService) FileService {
	return &fileService{
		repo:      repo,
		collaRepo: collaRepo,
		uSvc:      uSvc,
	}
}

func (s *fileService) GetByUserId(userId string) ([]datamodels.File, bool) {
	collaborators, found := s.collaRepo.GetByUserId(userId)
	if !found {
		return nil, false
	}

	var fileIds []uint
	for _, c := range collaborators {
		fileIds = append(fileIds, c.FileId)
	}

	files, found := s.repo.BatchGet(fileIds)
	if !found {
		return nil, false
	}

	return files, true
}

func (s *fileService) GetByFileId(fileId uint) (datamodels.File, bool) {
	return s.repo.Get(fileId)
}

func (s *fileService) Create(req CreateUnitRequest) (datamodels.File, error) {
	unitId, err := s.uSvc.CreateUnit(req)
	if err != nil {
		log.Printf("Error while creating unit: %v", err)
		return datamodels.File{}, err
	}

	file := datamodels.File{
		Name:     req.Name,
		UnitType: datamodels.FileTypeInt(req.Type),
		UnitId:   unitId,
	}

	file, err = s.repo.Create(file)
	if err != nil {
		log.Printf("Error while creating file: %v", err)
		return datamodels.File{}, err
	}

	_, err = s.collaRepo.Create(datamodels.FileCollaborator{
		FileId: file.ID,
		UserId: req.UserId,
		Role:   datamodels.RoleOwner,
	})
	if err != nil {
		log.Printf("Error while creating file collaborator: %v", err)
		return datamodels.File{}, err
	}

	return file, nil
}

func (s *fileService) GetCollaborators(fileId uint) ([]datamodels.FileCollaborator, bool) {
	return s.collaRepo.GetByFileId(fileId)
}

func (s *fileService) GetCollaboratorsByUnitId(unitId string) ([]datamodels.FileCollaborator, bool) {
	file, found := s.repo.GetByUnitId(unitId)
	if !found {
		return nil, false
	}
	return s.collaRepo.GetByFileId(file.ID)
}
