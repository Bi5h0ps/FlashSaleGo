package model

type Product struct {
	ID           int64  `json:"ID" gorm:"column:ID;primaryKey;autoIncrement" iRaiden:"ID"`
	ProductName  string `json:"ProductName" gorm:"column:productName" iRaiden:"ProductName"`
	ProductNum   int64  `json:"ProductNum" gorm:"column:productNum" iRaiden:"ProductNum"`
	ProductImage string `json:"ProductImage" gorm:"column:productImage" iRaiden:"ProductImage"`
	ProductUrl   string `json:"ProductUrl" gorm:"column:productUrl" iRaiden:"ProductUrl"`
}
