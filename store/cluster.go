package store

import "DVT_Service/models"

// CreateClusterAmountTask Add data in table ClusterAmountTask
func (s *Store) CreateClusterAmountTask(tasks *models.ClusterAmountTask) error {
	return s.mysql.Create(tasks).Error
}

// CloseClusterAmountTask Unlock ClusterAmountTask by primary key
func (s *Store) CloseClusterAmountTask(taskId, err string) error {
	return s.mysql.Model(&models.ClusterAmountTask{}).
		Where("uuid = ?", taskId).
		Update("errorMsg", err).
		Error
}

// FinishClusterAmountTask End ClusterAmountTask by primary key and update the txhash and operators, amount
func (s *Store) FinishClusterAmountTask(taskId, operators, txhash, amount string) error {
	return s.mysql.Model(&models.ClusterAmountTask{}).
		Where("uuid = ?", taskId).
		Updates(map[string]interface{}{"txhash": txhash, "operators": operators, "amount": amount}).
		Error
}
