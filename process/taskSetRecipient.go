package process

import (
	"context"
	"time"

	"github.com/bytedance/sonic"
	"github.com/ethereum/go-ethereum/common"
	"github.com/hibiken/asynq"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const Task_SetRecipient = "task:SetRecipient"

type TaskSetRecipient struct {
	UUID    string `json:"uuid"`
	Address string `json:"address"`
}

func NewSetRecipientTask(uuid, address string) (*asynq.Task, error) {
	payload, err := sonic.Marshal(&TaskSetRecipient{UUID: uuid, Address: address})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(Task_SetRecipient, payload), nil
}

func (p *Payload) HandleTaskSetRecipient(ctx context.Context, t *asynq.Task) error {
	ti := time.Now()
	var task TaskSetRecipient
	err := sonic.Unmarshal(t.Payload(), &task)
	if err != nil {
		log.Err(err).
			Str("type", "SetRecipient").
			Str("payload", string(t.Payload())).
			Msg("Unmarshal error")
		return nil
	}

	log.Info().
		Str("uuid", task.UUID).
		Str("address", task.Address).
		Str("type", "SetRecipient").
		Msg("task start")

	//p.ssvContract.SetFeeRecipientAddress(nil, common.HexToAddress(task.Address))
	input, _ := p.ssvAbi.Pack("setFeeRecipientAddress", common.HexToAddress(task.Address))
	tx, err := p.kms.SignAndSendTransactionByKMSCtx(ctx, p.ssvAddr, input)
	if err != nil {
		log.Err(err).Msg("doSetFeeRecipient KMSSignAndSendTransactionByKMSCtx error")
		p.store.CloseFeeRecipientTask(task.UUID, errors.Wrap(err, "Transaction error").Error())
		return nil
	}

	p.store.FinishFeeRecipientTask(task.UUID, task.Address, tx.TxHash.Hex())

	log.Info().
		Str("uuid", task.UUID).
		Str("txhash", tx.TxHash.Hex()).
		Str("type", "SetRecipient").
		Str("timeuse", time.Now().Sub(ti).String()).
		Msg("task success")

	return nil
}
