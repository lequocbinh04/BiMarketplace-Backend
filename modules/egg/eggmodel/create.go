package eggmodel

import (
	"BiMarketplace/appCommon"
	"github.com/shopspring/decimal"
)

const EntityName = "Egg"

type Egg struct {
	appCommon.SQLModel `json:",inline"`
	OwnerAddress       string              `json:"owner_address" gorm:"column:owner_address"`
	TxHash             string              `json:"tx_hash" gorm:"column:tx_hash"`
	BlockNumber        int64               `json:"block_number" gorm:"column:block_number"`
	PayableToken       string              `json:"payable_token" gorm:"column:payable_token"`
	Price              decimal.NullDecimal `json:"price" gorm:"column:price"`
}

func (Egg) TableName() string { return "eggs" }
