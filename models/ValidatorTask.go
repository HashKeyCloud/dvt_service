package models

import (
	"time"
)

// ClusterValidatorTask Entity class of table validator_tasks
// foreign key: ValidatorInfo's ID
type ClusterValidatorTask struct {
	UUID        string    `json:"uuid" gorm:"primaryKey;column:uuid;"`
	Type        uint      `json:"-" gorm:"index;not null;type:tinyint;"`
	ValidatorID uint      `json:"-" gorm:"column:validator_id;"`
	Operators   string    `json:"operators" gorm:"not null;column:operators;type:varchar(200);"`
	TxHash      string    `json:"txHash" gorm:"column:txhash;type:varchar(100);"`
	Error       string    `json:"error" gorm:"column:errorMsg;"`
	CreatedAt   time.Time `json:"-" `
	UpdatedAt   time.Time `json:"time" gorm:"index"`
}

// ValidatorTaskOutput Struct ValidatorTask output struct
type ValidatorTaskOutput struct {
	UUID      string    `json:"uuid"`
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	TxHash    string    `json:"txHash"`
	Error     string    `json:"errorMsg,omitempty"`
}

// Process Turn ValidatorTask into ValidatorTaskOutput
func (t *ClusterValidatorTask) Process() *ValidatorTaskOutput {
	var typeStr string

	switch t.Type {
	case 1:
		typeStr = "registerValidator"
	case 2:
		typeStr = "removeValidator"
	}

	return &ValidatorTaskOutput{
		UUID:      t.UUID,
		Type:      typeStr,
		TxHash:    t.TxHash,
		Error:     t.Error,
		Timestamp: t.UpdatedAt,
	}
}
