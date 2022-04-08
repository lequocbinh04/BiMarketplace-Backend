package petstorage

import (
	"BiMarketplace/appCommon"
	"BiMarketplace/modules/pet/petmodel"
	"context"
	"gorm.io/gorm"
)

func (s *sqlStore) GetDataWithCondition(ctx context.Context, condition map[string]interface{}, moreKeys ...string) (*petmodel.Pet, error) {
	db := s.db

	var pet petmodel.Pet

	for i := range moreKeys {
		db = db.Preload(moreKeys[i]) // for auto preload
	}

	if err := db.Where(condition).First(&pet).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, appCommon.RecordNotFound
		}
		return nil, appCommon.ErrDB(err)
	}
	return &pet, nil
}
