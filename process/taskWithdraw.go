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

const Task_Withdraw = "task:Withdraw"

type TaskWithdraw struct {
	UUID   string `json:"uuid"`
	Amount string `json:"amount"`
}

func NewWithdrawTask(uuid, amount string) (*asynq.Task, error) {
	payload, err := sonic.Marshal(&TaskWithdraw{UUID: uuid, Amount: amount})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(Task_Withdraw, payload), nil
}

func (p *Payload) HandleTaskWithdraw(ctx context.Context, t *asynq.Task) error {
	ti := time.Now()
	var task TaskWithdraw
	err := sonic.Unmarshal(t.Payload(), &task)
	if err != nil {
		log.Err(err).
			Str("type", "Withdraw").
			Str("payload", string(t.Payload())).
			Msg("Unmarshal error")
		return nil
	}

	withdrawAmount, b := new(big.Int).SetString(task.Amount, 10)
	if !b {
		p.store.CloseClusterAmountTask(task.UUID, "invalid amount")
		return nil
	}

	log.Info().
		Str("uuid", task.UUID).
		Str("amount", task.Amount).
		Str("type", "Withdraw").
		Msg("task start")

	cluster, err := p.getCluster(ctx, p.defaultOperatorStr)
	if err != nil {
		p.store.CloseClusterAmountTask(task.UUID, errors.Wrap(err, "scanCluster error").Error())
		return nil
	}

	if !cluster.Active {
		log.Warn().Any("cluster", cluster).Msg("doWithdraw skip: the cluster is not active")
		p.store.CloseClusterAmountTask(task.UUID, "the cluster is not active")
		return nil
	}

	//p.ssvContract.Withdraw(nil, p.defaultOperatorIds, withdrawAmount, *cluster)
	input, _ := p.ssvAbi.Pack("withdraw", p.defaultOperatorIds, withdrawAmount, *cluster)
	tx, err := p.kms.SignAndSendTransactionByKMSCtx(ctx, p.ssvAddr, input)
	if err != nil {
		log.Err(err).Msg("doWithdraw KMSSignAndSendTransactionByKMSCtx error")
		p.store.DelSSVClusterSnapshot(ctx, p.defaultOperatorStr)
		p.store.CloseClusterAmountTask(task.UUID, errors.Wrap(err, "Transaction error").Error())
		return nil
	}

	p.updateClusterSnapshot(ctx, tx)

	p.store.FinishClusterAmountTask(task.UUID, p.defaultOperatorStr, tx.TxHash.Hex(), withdrawAmount.String())

	log.Info().
		Str("uuid", task.UUID).
		Str("txhash", tx.TxHash.Hex()).
		Str("type", "Withdraw").
		Str("timeuse", time.Now().Sub(ti).String()).
		Msg("task success")

	return nil
}
