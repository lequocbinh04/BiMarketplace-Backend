package transactionstorage

import (
	"BiMarketplace/modules/transaction/transactionmodel"
)

func (s *sqlStore) CreateNewTx(data *transactionmodel.Transaction) error {
	db := s.db
	if err := db.Create(data).Error; err != nil {
		return err
	}
	return nil
}
