package store

import (
	"DVT_Service/models"
)

// CreateValidatorTasks Add data in table ClusterValidatorTask in batches
func (s *Store) CreateValidatorTasks(tasks []*models.ClusterValidatorTask) error {
	return s.mysql.Create(tasks).Error
}

// CloseValidatorTask Unlock ClusterValidatorTask based on primary key
func (s *Store) CloseValidatorTask(taskId, err string) error {
	return s.mysql.Model(&models.ClusterValidatorTask{}).
		Where("uuid = ?", taskId).
		Update("errorMsg", err).
		Error
}

// FinishValidatorTask End ClusterValidatorTask based on primary key and update the txhash
func (s *Store) FinishValidatorTask(taskId, operators, txhash string) error {
	return s.mysql.Model(&models.ClusterValidatorTask{}).
		Where("uuid = ?", taskId).
		Updates(map[string]interface{}{"txhash": txhash, "operators": operators}).
		Error
}
