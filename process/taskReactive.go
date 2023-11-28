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

const Task_Reactive = "task:Reactive"

type TaskReactive struct {
	UUID string `json:"uuid"`
}

func NewReactiveTask(uuid string) (*asynq.Task, error) {
	payload, err := sonic.Marshal(&TaskReactive{UUID: uuid})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(Task_Reactive, payload), nil
}

func (p *Payload) HandleTaskReactive(ctx context.Context, t *asynq.Task) error {
	ti := time.Now()
	var task TaskReactive
	err := sonic.Unmarshal(t.Payload(), &task)
	if err != nil {
		log.Err(err).
			Str("type", "Reactive").
			Str("payload", string(t.Payload())).
			Msg("Unmarshal error")
		return nil
	}

	log.Info().
		Str("uuid", task.UUID).
		Str("type", "Reactive").
		Msg("task start")

	p.store.DelSSVClusterSnapshot(ctx, p.defaultOperatorStr)
	cluster, err := p.getCluster(ctx, p.defaultOperatorStr)
	if err != nil {
		log.Err(err).Msg("scanCluster error")
		p.store.CloseClusterAmountTask(task.UUID, errors.Wrap(err, "scanCluster error").Error())
		return nil
	}

	if cluster.Active {
		log.Warn().Any("cluster", cluster).Msg("doReactive skip: the cluster is active")
		p.store.CloseClusterAmountTask(task.UUID, "the cluster is active")
		return nil
	}

	ValidatorCount := new(big.Int).SetInt64(int64(cluster.ValidatorCount))
	reactiveAmount := new(big.Int).Mul(p.amountTokenSSV, ValidatorCount)

	//p.ssvContract.Reactivate(nil, p.defaultOperatorIds, reactiveAmount, *cluster)
	input, _ := p.ssvAbi.Pack("reactivate", p.defaultOperatorIds, reactiveAmount, *cluster)
	tx, err := p.kms.SignAndSendTransactionByKMSCtx(ctx, p.ssvAddr, input)
	if err != nil {
		log.Err(err).Msg("doReactive KMSSignAndSendTransactionByKMSCtx error")
		p.store.DelSSVClusterSnapshot(ctx, p.defaultOperatorStr)
		p.store.CloseClusterAmountTask(task.UUID, errors.Wrap(err, "Transaction error").Error())
		return nil
	}
	p.updateClusterSnapshot(ctx, tx)

	p.store.FinishClusterAmountTask(task.UUID, p.defaultOperatorStr, tx.TxHash.Hex(), reactiveAmount.String())

	log.Info().
		Str("uuid", task.UUID).
		Str("txhash", tx.TxHash.Hex()).
		Str("type", "Reactive").
		Str("timeuse", time.Now().Sub(ti).String()).
		Msg("task success")

	return nil
}
