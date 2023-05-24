package service

import (
	"FlashSaleGo/model"
	"FlashSaleGo/repository"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	IsPwdSuccess(userName string, pwd string) (user *model.User, isOk bool)
	AddUser(user *model.User) (userId int64, err error)
	GetUserByUsername(username string) (*model.User, error)
}

type UserService struct {
	UserRepository repository.IUser
}

func (u *UserService) IsPwdSuccess(userName string, pwd string) (user *model.User, isOk bool) {
	var err error
	user, err = u.UserRepository.Select(userName)
	if err != nil {
		return
	}
	isOk, _ = ValidatePassword(pwd, user.HashPassword)
	if !isOk {
		return &model.User{}, false
	}
	return
}

func (u *UserService) AddUser(user *model.User) (userId int64, err error) {
	//user.HashPassword is user's input here, unhashed
	pwdByte, err := GeneratePassword(user.HashPassword)
	if err != nil {
		return 0, err
	}
	//replace password by hashed password and write to the database
	user.HashPassword = string(pwdByte)
	return u.UserRepository.Insert(user)
}

func (u *UserService) GetUserByUsername(username string) (user *model.User, err error) {
	return u.UserRepository.Select(username)
}

// ValidatePassword compares hashed password stored in the database with user provided password
func ValidatePassword(userPassword string, hashed string) (isOk bool, err error) {
	if err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(userPassword)); err != nil {
		return false, errors.New("password incorrect")
	}
	return true, nil
}

// GeneratePassword generates a hashed user password
func GeneratePassword(userPassword string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
}

func NewUserService(repository repository.IUser) IUserService {
	return &UserService{UserRepository: repository}
}
