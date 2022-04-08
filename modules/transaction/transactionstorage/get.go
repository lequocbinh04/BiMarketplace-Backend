package transactionstorage

import (
	"BiMarketplace/appCommon"
	"BiMarketplace/modules/transaction/transactionmodel"
	"context"
	"gorm.io/gorm"
)

func (s *sqlStore) GetDataWithCondition(ctx context.Context, condition map[string]interface{}, moreKeys ...string) (*transactionmodel.Transaction, error) {
	db := s.db

	var tx transactionmodel.Transaction

	for i := range moreKeys {
		db = db.Preload(moreKeys[i]) // for auto preload
	}

	if err := db.Where(condition).First(&tx).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, appCommon.RecordNotFound
		}
		return nil, appCommon.ErrDB(err)
	}
	return &tx, nil
}
