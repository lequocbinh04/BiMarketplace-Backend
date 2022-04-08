package eggstorage

import (
	"BiMarketplace/appCommon"
	"BiMarketplace/modules/egg/eggmodel"
	"context"
	"gorm.io/gorm"
)

func (s *sqlStore) GetDataWithCondition(ctx context.Context, condition map[string]interface{}, moreKeys ...string) (*eggmodel.Egg, error) {
	db := s.db

	var egg eggmodel.Egg

	for i := range moreKeys {
		db = db.Preload(moreKeys[i]) // for auto preload
	}

	if err := db.Where(condition).First(&egg).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, appCommon.RecordNotFound
		}
		return nil, appCommon.ErrDB(err)
	}
	return &egg, nil
}
