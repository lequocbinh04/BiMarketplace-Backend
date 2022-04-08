package userstorage

import (
	"BiMarketplace/appCommon"
	"BiMarketplace/modules/user/usermodel"
	"context"
	"gorm.io/gorm"
)

func (s *sqlStore) FindUser(ctx context.Context, Address string) (*usermodel.User, error) {
	db := s.db
	var user usermodel.User
	if err := db.Where("wallet_address = ?", Address).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, appCommon.RecordNotFound
		}
		return nil, appCommon.ErrDB(err)
	}
	return &user, nil
}
