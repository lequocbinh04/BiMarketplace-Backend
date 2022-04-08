package usermodel

import "BiMarketplace/appCommon"

const EntityName = "User"
const NonceEntityName = "Nonce"

type User struct {
	appCommon.SQLModel `json:",inline"`
	WalletAddress      string           `json:"wallet_address" gorm:"column:wallet_address;"`
	Nonce              int              `json:"nonce" gorm:"column:nonce;"`
	DisplayName        string           `json:"display_name" gorm:"column:display_name;"`
	Role               string           `json:"role" gorm:"column:role;"`
	Avatar             *appCommon.Image `json:"avatar,omitempty" gorm:"column:avatar;type:json"`
}

func (User) TableName() string { return "users" }

type UserLogin struct {
	Address   string `json:"address" gorm:"column:wallet_address" binding:"required"`
	Signature string `json:"signature" gorm:"column:signature" binding:"required"`
}

func (UserLogin) TableName() string { return "users" }
