package services

import (
	"errors"
	"image"

	"go-usip/datamodels"
	"go-usip/repositories"
)

// UserService handles CRUID operations of a user datamodel,
// it depends on a user repository for its actions.
// It's here to decouple the data source from the higher level compoments.
// As a result a different repository type can be used with the same logic without any aditional changes.
// It's an interface and it's used as interface everywhere
// because we may need to change or try an experimental different domain logic at the future.
type UserService interface {
	GetByID(userId string) (datamodels.User, bool)
	GetInIDs(userIds []string) ([]datamodels.User, bool)
	GetByUsernameAndPassword(username, userPassword string) (datamodels.User, bool)
	GetByPage(nextId, size uint) (users []datamodels.User, latest bool)
	GetAvatarByUserID(userId string) (image.Image, bool)

	DeleteByID(userId string) bool

	Update(userId string, user datamodels.User) (datamodels.User, error)
	UpdatePassword(userId string, newPassword string) (datamodels.User, error)
	UpdateUsername(userId string, newUsername string) (datamodels.User, error)

	Create(userPassword string, user datamodels.User) (datamodels.User, error)
}

// NewUserService returns the default user service.
func NewUserService(repo repositories.UserRepository, aSvc AvatarService) UserService {
	return &userService{
		repo: repo,
		aSvc: aSvc,
	}
}

type userService struct {
	repo repositories.UserRepository
	aSvc AvatarService
}

// GetByID returns a user based on its id.
func (s *userService) GetByID(userId string) (datamodels.User, bool) {
	return s.repo.Get(userId)
}

func (s *userService) GetInIDs(userIds []string) ([]datamodels.User, bool) {
	return s.repo.BatchGet(userIds)
}

// GetByUsernameAndPassword returns a user based on its username and password,
// used for authentication.
func (s *userService) GetByUsernameAndPassword(username, userPassword string) (datamodels.User, bool) {
	if username == "" || userPassword == "" {
		return datamodels.User{}, false
	}

	user, found := s.repo.GetByUsername(username)
	if !found {
		return datamodels.User{}, false
	}
	if ok, _ := datamodels.ValidatePassword(userPassword, user.HashedPassword); !ok {
		return datamodels.User{}, false
	}

	return user, true
}

func (s *userService) GetByPage(nextId, size uint) ([]datamodels.User, bool) {
	users, found := s.repo.GetByPage(nextId, size+1)
	if !found {
		return nil, false
	}
	if len(users) < int(size)+1 {
		return users, true
	}
	return users[:size], false
}

func (s *userService) GetAvatarByUserID(userId string) (image.Image, bool) {
	user, found := s.repo.Get(userId)
	if !found {
		return nil, false
	}

	name := []rune(user.Nickname)
	c := name[0:1]
	image, err := s.aSvc.GenerateAvatar(string(c))
	if err != nil {
		return nil, false
	}

	return image, true
}

// Update updates every field from an existing User,
// it's not safe to be used via public API,
// however we will use it on the web/controllers/user_controller.go#PutBy
// in order to show you how it works.
func (s *userService) Update(userId string, user datamodels.User) (datamodels.User, error) {
	user.UserId = userId
	return s.repo.InsertOrUpdate(user)
}

// UpdatePassword updates a user's password.
func (s *userService) UpdatePassword(userId string, newPassword string) (datamodels.User, error) {
	// update the user and return it.
	hashed, err := datamodels.GeneratePassword(newPassword)
	if err != nil {
		return datamodels.User{}, err
	}

	return s.Update(userId, datamodels.User{
		HashedPassword: hashed,
	})
}

// UpdateUsername updates a user's username.
func (s *userService) UpdateUsername(userId string, newUsername string) (datamodels.User, error) {
	return s.Update(userId, datamodels.User{
		Username: newUsername,
	})
}

// Create inserts a new User,
// the userPassword is the client-typed password
// it will be hashed before the insertion to our repository.
func (s *userService) Create(userPassword string, user datamodels.User) (datamodels.User, error) {
	if user.ID > 0 || userPassword == "" || user.Nickname == "" || user.Username == "" {
		return datamodels.User{}, errors.New("unable to create this user")
	}

	user.UserId = datamodels.GenerateUserId()

	hashed, err := datamodels.GeneratePassword(userPassword)
	if err != nil {
		return datamodels.User{}, err
	}
	user.HashedPassword = hashed

	return s.repo.InsertOrUpdate(user)
}

// DeleteByID deletes a user by its id.
//
// Returns true if deleted otherwise false.
func (s *userService) DeleteByID(userId string) bool {
	return s.repo.Delete(userId)
}
