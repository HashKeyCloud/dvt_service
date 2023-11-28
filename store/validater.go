package store

import (
	"DVT_Service/models"
)

// SaveValidator Construct transactions, add and update bulk operation table ValidatorInfo
func (s *Store) SaveValidator(keys, links, pubkeys, encrypted []string) error {
	keystores := make([]*models.ValidatorInfo, len(keys))
	linkInfos := make([]*models.Encrypt, len(keys))

	for i := range keys {
		keystores[i] = &models.ValidatorInfo{
			PublicKey: pubkeys[i],
			Keystore:  keys[i],
			Operators: "",
		}

		linkInfos[i] = &models.Encrypt{
			Link:              links[i],
			EncryptedPassword: encrypted[i],
		}
	}

	tx := s.mysql.Begin()

	if err := tx.Create(&keystores).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Create(&linkInfos).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (s *Store) GetEncryptedPassword(hash string) (string, error) {
	var res models.Encrypt
	if err := s.mysql.Where("link = ?", hash).First(&res).Error; err != nil {
		return "", err
	} else {
		return res.EncryptedPassword, nil
	}
}

// GetValidatorStateByPublicKey Use publicKey to get Validator info
func (s *Store) GetValidatorStateByPublicKey(publicKey string) (*models.ValidatorInfo, error) {
	var res *models.ValidatorInfo
	err := s.mysql.Model(&models.ValidatorInfo{}).Preload("Tasks").Where("publicKey = ?", publicKey).First(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Store) GetValidatorStateByID(id uint) ([]*models.ValidatorInfo, error) {
	var res []*models.ValidatorInfo
	err := s.mysql.Model(&models.ValidatorInfo{}).
		Where("id = ?", id).
		Find(&res).
		Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetValidatorInfoByPublicKey Use publicKeys to get pending Validator info
func (s *Store) GetValidatorInfoByPublicKey(publicKeys []string) ([]*models.ValidatorInfo, error) {
	var res []*models.ValidatorInfo
	err := s.mysql.Model(&models.ValidatorInfo{}).Where("publicKey in ?", publicKeys).Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

// PendingValidatorState Batch pending Validator info
func (s *Store) PendingValidatorState(keys []*models.ValidatorInfo) error {
	var ids []uint
	for _, key := range keys {
		ids = append(ids, key.ID)
	}

	err := s.mysql.Model(&models.ValidatorInfo{}).Where("id in ?", ids).
		UpdateColumn("state", 1).Error
	return err
}

// StakingValidatorState staking Validator info by primary key
func (s *Store) StakingValidatorState(id uint, operator string) error {
	err := s.mysql.Model(&models.ValidatorInfo{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"state":     2,
			"operators": operator,
		}).Error
	return err
}

// TurnBackValidatorState Rollback Validator info state
func (s *Store) TurnBackValidatorState(id, state uint) error {
	err := s.mysql.Model(&models.ValidatorInfo{}).Where("id = ?", id).
		UpdateColumn("state", state).Error
	return err
}

// BackValidatorState Rollback Validator info state
func (s *Store) BackValidatorState(id uint, ptime int64) error {
	err := s.mysql.Model(&models.ValidatorInfo{}).Where("id = ?", id).
		Updates(map[string]interface{}{"state": 0, "pendingTime": ptime}).Error
	return err
}
