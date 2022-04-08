package configstorage

import (
	"BiMarketplace/appCommon"
	"BiMarketplace/modules/config/configmodel"
	"context"
)

func (s *sqlStore) UpdateLatestBlockNumber(ctx context.Context, num int64) error {
	db := s.db
	if err := db.Table(configmodel.GetConfig{}.TableName()).Where("id = ?", 1).Update("latest_block_number", num).Error; err != nil {
		return appCommon.ErrDB(err)
	}
	return nil
}
