package usermodel

type Nonce struct {
	Address string `json:"address" gorm:"column:wallet_address;"`
	Nonce   int    `json:"nonce" gorm:"column:nonce;"`
}
