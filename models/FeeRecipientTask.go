package models

import "time"

// FeeRecipientTask Entity class of table fee_recipient_tasks
type FeeRecipientTask struct {
	UUID         string    `json:"uuid" gorm:"primaryKey;column:uuid;"`
	FeeRecipient string    `json:"feeRecipient" gorm:"column:fee_recipient;type:varchar(100);"`
	TxHash       string    `json:"txHash" gorm:"column:txhash;type:varchar(100);"`
	Error        string    `json:"error" gorm:"column:errorMsg;"`
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"time" gorm:"index"`
}
