package repositories

import (
	"go-usip/datamodels"
	"log"

	"gorm.io/gorm"
)

// UserRepository handles the basic operations of a user entity/model.
// It's an interface in order to be testable, i.e a memory user repository or
// a connected to an sql database.
type UserRepository interface {
	Get(userId string) (user datamodels.User, found bool)
	BatchGet(userIds []string) (users []datamodels.User, found bool)
	GetByUsername(username string) (user datamodels.User, found bool)

	InsertOrUpdate(user datamodels.User) (updatedUser datamodels.User, err error)
	Delete(userId string) (deleted bool)
}

// NewUserRepository returns a new user memory-based repository,
// the one and only repository type in our example.
func NewUserRepository(db *gorm.DB) UserRepository {
	if err := db.AutoMigrate(&datamodels.User{}); err != nil {
		log.Fatalf("AutoMigrate error: %v", err)
	}

	return &userRepository{db: db}
}

type userRepository struct {
	db *gorm.DB
}

func (r *userRepository) Get(userId string) (user datamodels.User, found bool) {
	user = datamodels.User{}
	if err := r.db.Where("user_id = ?", userId).First(&user).Error; err != nil {
		log.Printf("Error while getting user by id: %v", err)
		return user, false
	}
	return user, true
}

func (r *userRepository) BatchGet(userIds []string) (users []datamodels.User, found bool) {
	if err := r.db.Where("user_id IN ?", userIds).Find(&users).Error; err != nil {
		log.Printf("Error while getting users by ids: %v", err)
		return users, false
	}
	return users, true
}

func (r *userRepository) GetByUsername(username string) (user datamodels.User, found bool) {
	user = datamodels.User{}
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		log.Printf("Error while getting user by username: %v", err)
		return user, false
	}
	return user, true
}

func (r *userRepository) InsertOrUpdate(user datamodels.User) (datamodels.User, error) {
	if user.ID > 0 {
		return user, r.db.Save(&user).Error
	}

	return user, r.db.Create(&user).Error
}

func (r *userRepository) Delete(userId string) (deleted bool) {
	if err := r.db.Delete(&datamodels.User{}, userId).Error; err != nil {
		log.Printf("Error while deleting user by id: %v", err)
		return false
	}
	return true
}
