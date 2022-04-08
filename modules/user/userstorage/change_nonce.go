package userstorage

import (
	"BiMarketplace/appCommon"
	"BiMarketplace/modules/user/usermodel"
)

func (s *sqlStore) ChangeNonce(nonce int, userData *usermodel.User) error {
	db := s.db

	createData := userData
	createData.Nonce = nonce

	if err := db.Where("id = ?", userData.Id).Updates(createData).Error; err != nil {
		return appCommon.ErrDB(err)
	}

	return nil
}
