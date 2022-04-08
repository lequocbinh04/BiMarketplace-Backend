package eggstorage

import (
	"BiMarketplace/appCommon"
	"BiMarketplace/modules/egg/eggmodel"
)

func (s *sqlStore) Create(data *eggmodel.Egg) error {
	db := s.db
	if err := db.Create(data).Error; err != nil {
		return appCommon.ErrDB(err)
	}
	return nil
}
