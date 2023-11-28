package models

import "time"

// ClusterAmountTask Entity class of table fee_recipient_tasks
type ClusterAmountTask struct {
	UUID      string `gorm:"primaryKey;column:uuid;"`
	Type      uint   `gorm:"index;not null;type:tinyint;"`
	Operators string `gorm:"not null;column:operators;type:varchar(200);"`
	Amount    string `gorm:"not null;column:amount;type:varchar(100);"`
	TxHash    string `gorm:"column:txhash;type:varchar(100);"`
	Error     string `gorm:"column:errorMsg;"`
	CreatedAt time.Time
	UpdatedAt time.Time `gorm:"index"`
}
