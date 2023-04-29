package controller

import (
	"FlashSaleGo/model"
	"FlashSaleGo/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"os"
	"time"
)

type UserController struct {
	UserService service.IUserService
}

func (u *UserController) GetRegister(c *gin.Context) {
	c.HTML(http.StatusOK, "register", gin.H{})
}

func (u *UserController) PostRegister(c *gin.Context) {
	nickName := c.PostForm("nickName")
	userName := c.PostForm("userName")
	password := c.PostForm("password")

	user := &model.User{
		UserName:     userName,
		NickName:     nickName,
		HashPassword: password,
	}
	_, err := u.UserService.AddUser(user)
	if err != nil {
		c.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/user/login")
	return
}

func (u *UserController) GetLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login", gin.H{})
}

func (u *UserController) PostLogin(c *gin.Context) {
	//1. Extract user input
	var (
		userName = c.PostForm("userName")
		password = c.PostForm("password")
	)
	//2. Check if password is valid
	user, isOk := u.UserService.IsPwdSuccess(userName, password)
	if !isOk {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "Username or password not correct",
		})
	}
	//3.Generate jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.UserName,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	secretKeyByte := []byte(os.Getenv("SECRET"))
	tokenString, err := token.SignedString(secretKeyByte)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": "Failed to generate a signed token",
		})
		return
	}
	//4. Write jwt token to cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600, "", "", false, true)
	c.Redirect(http.StatusFound, "/product")
}
