package models

import (
	"sort"
	"time"

	"github.com/bytedance/sonic"
)

// ValidatorInfo Entity class of table validator_infos
type ValidatorInfo struct {
	ID          uint                    `gorm:"primaryKey;autoIncrement;column:id;"`
	PublicKey   string                  `gorm:"uniqueIndex;not null;column:publicKey;type:varchar(100);"`
	Keystore    string                  `gorm:"column:keystore;not null"`
	Operators   string                  `gorm:"column:operators;not null"`
	State       int                     `gorm:"index;not null;default:0;column:state;type:tinyint;"`
	PendingTime int64                   `gorm:"not null;default:0;column:pendingTime;"`
	Tasks       []*ClusterValidatorTask `gorm:"foreignKey:ValidatorID;references:ID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time `gorm:"index"`
}

type Encrypt struct {
	Link              string `json:"link" gorm:"column:link;primaryKey;type:varchar(100); not null"`
	EncryptedPassword string `json:"encrypted_password" gorm:"column:encrypted_password;type:text;not null"`
}

// ValidatorOutput Struct ValidatorInfo output struct
type ValidatorOutput struct {
	PublicKey   string                 `json:"public_key"`   // public string
	Operators   []int                  `json:"operators"`    // Operators
	State       string                 `json:"state"`        // state
	PendingTime int64                  `json:"pending_time"` // pending timestamp
	Tasks       []*ValidatorTaskOutput `json:"tasks"`        // tasks limit 5
}

// Process Turn ValidatorInfo into ValidatorOutput
func (v *ValidatorInfo) Process() *ValidatorOutput {
	sort.Slice(v.Tasks, func(i, j int) bool {
		return v.Tasks[i].UpdatedAt.Unix() > v.Tasks[j].UpdatedAt.Unix()
	})

	var tasks []*ValidatorTaskOutput
	for i, task := range v.Tasks {
		tasks = append(tasks, task.Process())
		if i == 4 {
			break
		}
	}

	var OperatorsArr []int
	sonic.UnmarshalString(v.Operators, &OperatorsArr)

	var StateStr string
	switch v.State {
	case 0:
		StateStr = "Unused"
	case 1:
		StateStr = "Processing"
	case 2:
		StateStr = "Staked"
	}

	return &ValidatorOutput{
		PublicKey:   v.PublicKey,
		Operators:   OperatorsArr,
		State:       StateStr,
		PendingTime: v.PendingTime,
		Tasks:       tasks,
	}
}
