package userstorage

import (
	"BiMarketplace/appCommon"
	"BiMarketplace/modules/user/usermodel"
	"gorm.io/gorm"
)

func (s *sqlStore) GetNonceDB(address string) (*usermodel.Nonce, error) {
	db := s.db
	var nonceDB usermodel.Nonce
	if err := db.Table(usermodel.User{}.TableName()).Where("wallet_address = ?", address).First(&nonceDB).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, appCommon.RecordNotFound
		}
		return nil, appCommon.ErrDB(err)
	}
	return &nonceDB, nil
}
