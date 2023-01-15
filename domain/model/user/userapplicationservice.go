package user

import (
	"fmt"
	"log"
	// "github.com/google/uuid"
)

type UserApplicationService struct {
	userRepository UserRepositorier
	userService    UserService
}

func NewUserApplicationService(userRepository UserRepositorier, userService UserService) (*UserApplicationService, error) {
	return &UserApplicationService{userRepository: userRepository, userService: userService}, nil
}

func (uas *UserApplicationService) Register(name string) (err error) {
	defer func() {
		if err != nil {
			err = &RegisterError{Name: name, Message: fmt.Sprintf("userapplicationservice.Register err: %s", err), Err: err}
		}
	}()
	userName, err := NewUserName(name)
	if err != nil {
		return err
	}

	// uuidV4 := uuid.New().String()
	uuidV4 := "test-id"
	userId, err := NewUserId(uuidV4)
	if err != nil {
		return err
	}

	user, err := NewUser(*userId, *userName)
	if err != nil {
		return err
	}

	isUserExists, err := uas.userService.Exists(user)
	if err != nil {
		return err
	}
	if isUserExists {
		return fmt.Errorf("user name of %s is already exists", name)
	}

	if err := uas.userRepository.Save(user); err != nil {
		return err
	}

	log.Printf("user name of %s is successfully saved", name)
	return nil
}

type RegisterError struct {
	Name    string
	Message string
	Err     error
}

func (err *RegisterError) Error() string {
	return err.Message
}

type UserData struct {
	Id   string
	Name string
}

func (uas *UserApplicationService) Get(userId string) (_ *UserData, err error) {
	defer func() {
		if err != nil {
			err = &GetError{UserId: userId, Message: fmt.Sprintf("userapplicationservice.Get err: %s", err), Err: err}
		}
	}()
	targetId, err := NewUserId(userId)
	if err != nil {
		return nil, err
	}
	user, err := uas.userRepository.FindByUserId(targetId)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, nil
	}
	userData := &UserData{Id: user.id.value, Name: user.name.value}
	return userData, nil
}

type GetError struct {
	UserId  string
	Message string
	Err     error
}

func (err *GetError) Error() string {
	return err.Message
}

type UserUpdateCommand struct {
	Id   string
	Name string
}

func (uas *UserApplicationService) Update(command UserUpdateCommand) error {
	targetId, err := NewUserId(command.Id)
	if err != nil {
		return err
	}
	user, err := uas.userRepository.FindByUserId(targetId)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("user is not found")
	}

	if name := command.Name; name != "" {
		newUserName, err := NewUserName(name)
		if err != nil {
			return err
		}
		user.ChangeName(*newUserName)

		isExists, err := uas.userService.Exists(user)
		if err != nil {
			return err
		}
		if isExists {
			return fmt.Errorf("user name of %s is already exists", name)
		}
	}

	if err := uas.userRepository.Update(user); err != nil {
		return err
	}

	log.Println("successfully updated")
	return nil
}

type UserDeleteCommand struct {
	Id string
}

func (uas *UserApplicationService) Delete(command UserDeleteCommand) error {
	targetId, err := NewUserId(command.Id)
	if err != nil {
		return err
	}
	user, err := uas.userRepository.FindByUserId(targetId)
	if err != nil {
		return err
	}
	if user == nil {
		return nil
	}

	if err := uas.userRepository.Delete(user); err != nil {
		return err
	}
	log.Println("successfully deleted")
	return nil
}
