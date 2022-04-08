package userstorage

import (
	"BiMarketplace/appCommon"
	"BiMarketplace/modules/user/usermodel"
)

func (s *sqlStore) CreateNonce(nonce int, address string) error {
	db := s.db

	createData := &usermodel.User{
		WalletAddress: address,
		DisplayName:   "",
		Role:          "USER",
		Nonce:         nonce,
		Avatar: &appCommon.Image{
			Url:    "",
			Width:  0,
			Height: 0,
		},
	}

	createData.Status = appCommon.UserStatusActivated

	if err := db.Create(createData).Error; err != nil {
		return appCommon.ErrDB(err)
	}

	return nil
}
