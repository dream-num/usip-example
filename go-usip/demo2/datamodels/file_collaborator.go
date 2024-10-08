package datamodels

type Role string

const (
	RoleOwner  Role = "owner"
	RoleEditor Role = "editor"
	RoleReader Role = "reader"
)

var RoleLever = map[Role]int{
	RoleOwner:  3,
	RoleEditor: 2,
	RoleReader: 1,
}

type FileCollaborator struct {
	ID     int64  `json:"id" gorm:"primary_key"`
	UserId string `json:"user_id" gorm:"uniqueIndex:uqe_file_id_user_id,piroity:2"`
	FileId uint   `json:"file_id" gorm:"uniqueIndex:uqe_file_id_user_id,piroity:1;index"`
	Role   Role   `json:"role"`
}
