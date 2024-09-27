package services

import (
	"errors"
	"go-usip/datamodels"
	"go-usip/repositories"
	"log"
	"mime/multipart"
	"strings"
)

type FileService interface {
	GetByUserId(userId string) ([]datamodels.File, bool)
	GetByFileId(fileId uint) (datamodels.File, bool)
	GetCollaborators(fileId uint) ([]datamodels.FileCollaborator, bool)
	GetCollaboratorsByUnitId(unitId string) ([]datamodels.FileCollaborator, bool)

	Create(req CreateUnitRequest) (datamodels.File, error)
	Import(req ImportReq) (datamodels.File, error)
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

func (s *fileService) create(unitId string, req CreateUnitRequest) (datamodels.File, error) {
	file := datamodels.File{
		Name:     req.Name,
		UnitType: datamodels.FileTypeInt(req.Type),
		UnitId:   unitId,
	}

	var err error
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

func (s *fileService) Create(req CreateUnitRequest) (datamodels.File, error) {
	unitId, err := s.uSvc.CreateUnit(req)
	if err != nil {
		log.Printf("Error while creating unit: %v", err)
		return datamodels.File{}, err
	}

	return s.create(unitId, req)
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

type ImportReq struct {
	FileName string
	FileSize int
	UserId   string
	Type     int

	FormFile multipart.File
	Cookie   string
}

func (s *fileService) Import(req ImportReq) (file datamodels.File, err error) {
	fileId, err := s.uSvc.UploadFile(req)
	if err != nil {
		log.Printf("Error while uploading file: %v", err)
		return
	}

	if fileId == "" {
		return file, errors.New("File upload failed, fileId is empty")
	}

	taskId, err := s.uSvc.Import(UniverserImportReq{
		FileId:     fileId,
		Type:       req.Type,
		OutputType: 1,
		Cookie:     req.Cookie,
	})

	if err != nil {
		log.Printf("Error while importing file: %v", err)
		return
	}

	if taskId == "" {
		return file, errors.New("File import failed, taskId is empty")
	}

	var unitId string
	for {
		unitId, err = s.uSvc.PullResult(UniverserPullReq{
			TaskId: taskId,
			Cookie: req.Cookie,
		})
		if err != nil {
			log.Printf("Error while getting task: %v", err)
			return
		}
		if unitId != "" {
			break
		}
	}

	file, err = s.create(unitId, CreateUnitRequest{
		Name:   strings.Split(req.FileName, ".")[0],
		Type:   datamodels.FileTypeStr(req.Type),
		UserId: req.UserId,
	})
	if err != nil {
		log.Printf("Error while creating file: %v", err)
		return
	}
	log.Printf("File created: %v", file)

	return file, nil
}
