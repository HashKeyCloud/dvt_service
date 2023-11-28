package store

import "DVT_Service/models"

// GetFeeRecipientTask Read a feeRecipientTask message with processing
func (s *Store) GetFeeRecipientTask() ([]*models.FeeRecipientTask, error) {
	var res []*models.FeeRecipientTask
	err := s.mysql.Model(&models.FeeRecipientTask{}).Where("state = 0").Order("id").Limit(1).Find(&res).Error
	return res, err
}

// CreateFeeRecipientTask Add data in table FeeRecipientTask in batches
func (s *Store) CreateFeeRecipientTask(tasks *models.FeeRecipientTask) error {
	return s.mysql.Create(tasks).Error
}

// SearchFeeRecipientTaskById Query FeeRecipientTask info by primary key
func (s *Store) SearchFeeRecipientTaskById(id string) ([]*models.FeeRecipientTask, error) {
	var res []*models.FeeRecipientTask
	err := s.mysql.Model(&models.FeeRecipientTask{}).Where("uuid = ?", id).Find(&res, id).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

// CloseFeeRecipientTask Unlock FeeRecipientTask by primary key
func (s *Store) CloseFeeRecipientTask(taskId, err string) error {
	return s.mysql.Model(&models.FeeRecipientTask{}).
		Where("uuid = ?", taskId).
		Update("errorMsg", err).
		Error
}

// FinishFeeRecipientTask End FeeRecipientTask by primary key and update the txhash
func (s *Store) FinishFeeRecipientTask(taskId, address, txhash string) error {
	return s.mysql.Model(&models.FeeRecipientTask{}).Where("uuid = ?", taskId).
		Updates(map[string]interface{}{"txhash": txhash, "fee_recipient": address}).Error
}
