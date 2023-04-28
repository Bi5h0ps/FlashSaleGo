package model

type User struct {
	ID           int64  `json:"id" form:"ID" gorm:"primaryKey;autoIncrement"`
	NickName     string `json:"nickName" form:"nickName" gorm:"column:nickName"`
	UserName     string `json:"userName" form:"userName" gorm:"primaryKey;column:userName"`
	HashPassword string `json:"-" form:"passWord" gorm:"column:passWord"`
}
