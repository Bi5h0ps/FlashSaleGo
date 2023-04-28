package model

type Order struct {
	ID          int64 `gorm:"primaryKey;autoIncrement"`
	UserId      int64 `gorm:"column:userId"`
	ProductId   int64 `gorm:"column:productId"`
	OrderStatus int   `gorm:"column:orderStatus"`
}

const (
	OrderWait = iota
	OrderSuccess
	OrderFailed
)
