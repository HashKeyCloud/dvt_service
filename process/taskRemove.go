package process

import (
	"context"
	"time"

	"github.com/bytedance/sonic"
	"github.com/ethereum/go-ethereum/common"
	"github.com/hibiken/asynq"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"DVT_Service/models"
)

const Task_Remove = "task:Remove"

type TaskRemove struct {
	UUID        string `json:"uuid"`
	ValidatorID uint   `json:"validator_id"`
}

func NewRemoveTask(uuid string, validatorId uint) (*asynq.Task, error) {
	payload, err := sonic.Marshal(&TaskRemove{UUID: uuid, ValidatorID: validatorId})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(Task_Remove, payload), nil
}

func (p *Payload) HandleTaskRemove(ctx context.Context, t *asynq.Task) error {
	ti := time.Now()
	var task TaskRemove
	err := sonic.Unmarshal(t.Payload(), &task)
	if err != nil {
		log.Err(err).
			Str("type", "Remove").
			Str("payload", string(t.Payload())).
			Msg("Unmarshal error")
		return nil
	}

	log.Info().
		Str("uuid", task.UUID).
		Str("type", "Remove").
		Uint("ValidatorID", task.ValidatorID).
		Msg("task start")

	var keystore *models.ValidatorInfo
	if keystores, err := p.store.GetValidatorStateByID(task.ValidatorID); err != nil {
		log.Err(err).Msg("GetValidatorStateByID error")
		p.store.CloseValidatorTask(task.UUID, errors.Wrap(err, "GetValidatorStateByID error").Error())
		p.store.TurnBackValidatorState(task.ValidatorID, 2)
		return nil
	} else if len(keystores) != 1 {
		log.Error().
			Uint("ValidatorID", task.ValidatorID).
			Msg("Validator not found")
		p.store.CloseValidatorTask(task.UUID, errors.Wrap(err, "Validator not found").Error())
		p.store.TurnBackValidatorState(task.ValidatorID, 2)
		return nil
	} else {
		keystore = keystores[0]
	}

	cluster, err := p.getCluster(ctx, keystore.Operators)
	if err != nil {
		log.Err(err).Msg("scanCluster error")
		p.store.CloseValidatorTask(task.UUID, errors.Wrap(err, "scanCluster error").Error())
		p.store.TurnBackValidatorState(task.ValidatorID, 2)
		return nil
	}

	if !cluster.Active {
		log.Warn().Any("cluster", cluster).Msg("doReactive skip: the cluster is not active")
		p.store.CloseValidatorTask(task.UUID, "the cluster is not active")
		p.store.TurnBackValidatorState(task.ValidatorID, 2)
		return nil
	}

	//p.ssvContract.RemoveValidator(nil, common.FromHex(keystore.PublicKey), p.defaultOperatorIds, *cluster)
	input, _ := p.ssvAbi.Pack("removeValidator", common.FromHex(keystore.PublicKey), p.defaultOperatorIds, *cluster)
	tx, err := p.kms.SignAndSendTransactionByKMSCtx(ctx, p.ssvAddr, input)
	if err != nil {
		log.Err(err).Msg("doRemove KMSSignAndSendTransactionByKMSCtx error")
		p.store.DelSSVClusterSnapshot(ctx, p.defaultOperatorStr)
		p.store.CloseValidatorTask(task.UUID, errors.Wrap(err, "Transaction error").Error())
		p.store.TurnBackValidatorState(task.ValidatorID, 2)
		return nil
	}
	p.updateClusterSnapshot(ctx, tx)

	p.store.FinishValidatorTask(task.UUID, p.defaultOperatorStr, tx.TxHash.Hex())
	newTaskTime := time.Now().UTC().Add(15 * time.Minute).Unix() // after 2 epoch time
	p.store.BackValidatorState(task.ValidatorID, newTaskTime)

	log.Info().
		Str("uuid", task.UUID).
		Str("txhash", tx.TxHash.Hex()).
		Str("type", "Remove").
		Str("timeuse", time.Now().Sub(ti).String()).
		Msg("task success")

	return nil
}
