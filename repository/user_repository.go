package repository

import (
	"FlashSaleGo/model"
	"errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type IUser interface {
	Conn() error
	Select(userName string) (user *model.User, err error)
	Insert(user *model.User) (userId int64, err error)
}

const defaultUserDB = "root:Nmdhj2e2d@tcp(127.0.0.1:3306)/flashSaleDB?charset=utf8"

type UserRepository struct {
	myGormConn *gorm.DB
}

func (p *UserRepository) Conn() (err error) {
	if p.myGormConn == nil {
		p.myGormConn, err = gorm.Open(mysql.Open(defaultUserDB), &gorm.Config{})
		if err != nil {
			return
		}
		err = p.myGormConn.AutoMigrate(&model.User{})
		if err != nil {
			return
		}
	}
	return nil
}

func (p *UserRepository) Select(userName string) (user *model.User, err error) {
	if err = p.Conn(); err != nil {
		return
	}
	user = &model.User{}
	if result := p.myGormConn.Where("userName", userName).First(user); result.Error != nil {
		return nil, err
	}
	return
}

func (p *UserRepository) Insert(user *model.User) (userId int64, err error) {
	if err = p.Conn(); err != nil {
		return
	}
	//not allowed to set user id
	user.ID = 0
	//check if user already exists
	checkUser, err := p.Select(user.UserName)
	if err != nil {
		return
	}
	if checkUser != nil {
		//user already exist
		return 0, errors.New("user already exist")
	}
	if result := p.myGormConn.Create(user); result.Error != nil {
		return 0, result.Error
	}
	return user.ID, err
}

func NewUserRepository(db *gorm.DB) IUser {
	return &UserRepository{myGormConn: db}
}
