package eggstorage

import "BiMarketplace/modules/egg/eggmodel"

func (s *sqlStore) DecreaseEgg(id int) error {
	db := s.db
	if err := db.Table(eggmodel.Egg{}.TableName()).
		Where("id", id).
		Update("status", "opened").Error; err != nil {
		return err
	}
	return nil
}
