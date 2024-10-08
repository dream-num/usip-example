package services

import (
	"errors"
	"go-usip/datamodels"
	"go-usip/repositories"
	"io"
	"log"
	"mime/multipart"
	"strings"
	"time"
)

type FileService interface {
	GetByUserId(userId string) ([]datamodels.File, bool)
	GetByFileId(fileId uint) (datamodels.File, bool)
	GetCollaborators(fileId uint) ([]datamodels.FileCollaborator, bool)
	GetCollaboratorsByUnitId(unitId string) ([]datamodels.FileCollaborator, bool)
	CheckPermission(req CheckPermissionReq) bool

	Create(req CreateUnitRequest) (datamodels.File, error)
	Import(req ImportReq) (datamodels.File, error)
	Export(req ExportReq) (resp ExportResp, err error)
	Join(req JoinReq) error

	BatchDelete(userId string, fileIds []uint) error
}

type fileService struct {
	repo      repositories.FileRepository
	collaRepo repositories.FileCollaboratorRepository

	uSvc UniverserService
}

func NewFileService(repo repositories.FileRepository, collaRepo repositories.FileCollaboratorRepository, uSvc UniverserService) FileService {
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
		time.Sleep(500 * time.Millisecond)
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

func (s *fileService) BatchDelete(userId string, fileIds []uint) error {
	return s.collaRepo.BatchDelete(userId, fileIds)
}

type ExportReq struct {
	FileId uint
	UserId string

	Cookie string
}

type ExportResp struct {
	FileName string
	Reader   io.ReadCloser
}

func (s *fileService) Export(req ExportReq) (resp ExportResp, err error) {
	file, found := s.GetByFileId(req.FileId)
	if !found {
		return resp, errors.New("File not found")
	}

	_, found = s.collaRepo.Get(req.FileId, req.UserId)
	if !found {
		return resp, errors.New("File not found")
	}

	taskId, err := s.uSvc.Export(UniverserExportReq{
		UnitId: file.UnitId,
		Type:   file.UnitType,
		Cookie: req.Cookie,
	})
	if err != nil {
		log.Printf("Error while exporting file: %v", err)
		return
	}

	var fileId string
	for {
		fileId, err = s.uSvc.PullResult(UniverserPullReq{
			TaskId:       taskId,
			Cookie:       req.Cookie,
			ExchangeType: ExchangeTypeExport,
		})
		if err != nil {
			log.Printf("Error while getting task: %v", err)
			return
		}
		if fileId != "" {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}

	reader, err := s.uSvc.GetFile(UniverserGetFileReq{
		FileId: fileId,
		Cookie: req.Cookie,
	})
	if err != nil {
		log.Printf("Error while getting file: %v", err)
		return
	}

	resp.Reader = reader
	resp.FileName = file.Name
	switch file.UnitType {
	case datamodels.UnitTypeDoc:
		resp.FileName += ".docx"
	case datamodels.UnitTypeSheet:
		resp.FileName += ".xlsx"
	}
	return
}

type JoinReq struct {
	UserIds []string
	FileId  uint
	Role    datamodels.Role
}

func (s *fileService) Join(req JoinReq) error {
	var data []datamodels.FileCollaborator
	for _, userId := range req.UserIds {
		data = append(data, datamodels.FileCollaborator{
			FileId: req.FileId,
			UserId: userId,
			Role:   req.Role,
		})
	}

	return s.collaRepo.InsertOrUpdate(data)
}

type Action string

const (
	ActionDelete Action = "delete"
	ActionJoin   Action = "join"
)

type CheckPermissionReq struct {
	FileId uint
	UserId string
	Action Action
}

func (s *fileService) CheckPermission(req CheckPermissionReq) bool {
	colla, found := s.collaRepo.Get(req.FileId, req.UserId)
	if !found {
		return false
	}

	switch req.Action {
	case ActionDelete:
		return colla.Role == datamodels.RoleOwner
	case ActionJoin:
		return datamodels.RoleLever[colla.Role] >= datamodels.RoleLever[datamodels.RoleEditor]
	}

	return false
}
