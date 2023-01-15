package user

type UserService struct {
	userRepository UserRepositorier
}

func NewUserService(userRepository UserRepositorier) (*UserService, error) {
	return &UserService{userRepository: userRepository}, nil
}

func (us *UserService) Exists(user *User) (bool, error) {
	user, err := us.userRepository.FindByUserName(&user.name)
	if err != nil {
		return false, err
	}
	return user != nil, nil
}
