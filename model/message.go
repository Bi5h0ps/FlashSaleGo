package model

type Message struct {
	ProductID int64
	UserID    int64
}

func NewMessage(userID, productID int64) *Message {
	return &Message{ProductID: productID, UserID: userID}
}
