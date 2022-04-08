package userstorage

import "BiMarketplace/modules/user/usermodel"

func (s *sqlStore) GetUserByAddress(address string) (*usermodel.User, error) {
	db := s.db
	var user usermodel.User
	if err := db.Where("wallet_address = ?", address).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
