package configstorage

import (
	"BiMarketplace/modules/config/configmodel"
	"context"
)

func (s *sqlStore) GetLatestBlockNumber(ctx context.Context) (int64, error) {
	db := s.db
	var config configmodel.GetConfig
	if err := db.Where("id = ?", 1).Find(&config).Error; err != nil {
		return 0, err
	}
	return config.LatestBlockNumber, nil
}
