package petstorage

import "BiMarketplace/modules/pet/petmodel"

func (s *sqlStore) CreateNewPet(data *petmodel.Pet) error {
	db := s.db
	if err := db.Create(data).Error; err != nil {
		return err
	}
	return nil
}
