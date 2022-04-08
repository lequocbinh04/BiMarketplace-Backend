package petmodel

import "BiMarketplace/appCommon"

const EntityName = "Pet"

type Pet struct {
	appCommon.SQLModel `json:",inline"`
	OwnerId            int64  `json:"ownerId" gorm:"column:owner_id"`
	MkpId              int64  `json:"mkpId" gorm:"column:mkp_id"`
	Image              string `json:"image" gorm:"column:image"`
	Name               string `json:"name" gorm:"column:name"`
	Hp                 int64  `json:"hp" gorm:"column:hp"`
	Attack             int64  `json:"attack" gorm:"column:attack"`
	Speed              int64  `json:"speed" gorm:"column:speed"`
	NftId              int64  `json:"nftId" gorm:"column:nft_id"`
}

func (Pet) TableName() string {
	return "pets"
}
