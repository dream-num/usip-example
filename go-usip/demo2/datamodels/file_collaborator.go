package datamodels

type Role string

const (
	RoleOwner  Role = "owner"
	RoleEditor Role = "editor"
	RoleReader Role = "reader"
)

type FileCollaborator struct {
	ID     int64  `json:"id" gorm:"primary_key"`
	UserId string `json:"user_id"`
	FileId uint   `json:"file_id"`
	Role   Role   `json:"role"`
}
