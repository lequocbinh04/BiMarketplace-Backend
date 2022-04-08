package transactionmodel

import "BiMarketplace/appCommon"

type Transaction struct {
	appCommon.SQLModel `json:",inline"`
	PetId              int64  `json:"pet_id" gorm:"column:pet_id"`
	OwnerId            int64  `json:"owner_id" gorm:"column:owner_id"`
	Price              int64  `json:"price" gorm:"column:price"`
	TxHash             string `json:"tx_hash" gorm:"column:tx_hash"`
	PayableToken       string `json:"payable_token" gorm:"column:payable_token"`
	BlockNumber        int64  `json:"block_number" gorm:"column:block_number"`
	SellerAddress      string `json:"seller_address" gorm:"column:seller_address"`
	BuyerAddress       string `json:"buyer_address" gorm:"column:buyer_address"`
}

func (Transaction) TableName() string {
	return "transactions"
}
