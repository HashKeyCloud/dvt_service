package process

import (
	"context"
	"math/big"
	"time"

	"github.com/bytedance/sonic"
	"github.com/hibiken/asynq"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const Task_Deposit = "task:Deposit"

type TaskDeposit struct {
	UUID   string `json:"uuid"`
	Amount string `json:"amount"`
}

func NewDepositTask(uuid, amount string) (*asynq.Task, error) {
	payload, err := sonic.Marshal(&TaskDeposit{UUID: uuid, Amount: amount})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(Task_Deposit, payload), nil
}

func (p *Payload) HandleTaskDeposit(ctx context.Context, t *asynq.Task) error {
	ti := time.Now()
	var task TaskDeposit
	err := sonic.Unmarshal(t.Payload(), &task)
	if err != nil {
		log.Err(err).
			Str("type", "Deposit").
			Str("payload", string(t.Payload())).
			Msg("Unmarshal error")
		return nil
	}

	depositAmount, b := new(big.Int).SetString(task.Amount, 10)
	if !b {
		p.store.CloseClusterAmountTask(task.UUID, "invalid amount")
		return nil
	}

	log.Info().
		Str("uuid", task.UUID).
		Str("amount", task.Amount).
		Str("type", "Deposit").
		Msg("task start")

	cluster, err := p.getCluster(ctx, p.defaultOperatorStr)
	if err != nil {
		p.store.CloseClusterAmountTask(task.UUID, errors.Wrap(err, "scanCluster error").Error())
		return nil
	}

	if !cluster.Active {
		log.Warn().Any("cluster", cluster).Msg("doDeposit skip: the cluster is not active")
		p.store.CloseClusterAmountTask(task.UUID, "the cluster is not active")
		return nil
	}

	//p.ssvContract.Deposit(nil, p.kms.PubKey, p.defaultOperatorIds, depositAmount, *cluster)
	input, _ := p.ssvAbi.Pack("deposit", p.kms.PubKey, p.defaultOperatorIds, depositAmount, *cluster)
	tx, err := p.kms.SignAndSendTransactionByKMSCtx(ctx, p.ssvAddr, input)
	if err != nil {
		log.Err(err).Msg("doDeposit KMSSignAndSendTransactionByKMSCtx error")
		p.store.DelSSVClusterSnapshot(ctx, p.defaultOperatorStr)
		p.store.CloseClusterAmountTask(task.UUID, errors.Wrap(err, "Transaction error").Error())
		return nil
	}

	p.updateClusterSnapshot(ctx, tx)

	p.store.FinishClusterAmountTask(task.UUID, p.defaultOperatorStr, tx.TxHash.Hex(), depositAmount.String())

	log.Info().
		Str("uuid", task.UUID).
		Str("txhash", tx.TxHash.Hex()).
		Str("type", "Deposit").
		Str("timeuse", time.Now().Sub(ti).String()).
		Msg("task success")

	return nil
}
