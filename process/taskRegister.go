package process

import (
	"context"
	"os/exec"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/hibiken/asynq"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"DVT_Service/models"
)

const Task_register = "task:Register"

type TaskRegister struct {
	UUID        string `json:"uuid"`
	ValidatorID uint   `json:"validator_id"`
}

func NewRegisterTask(uuid string, validatorId uint) (*asynq.Task, error) {
	payload, err := sonic.Marshal(&TaskRegister{UUID: uuid, ValidatorID: validatorId})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(Task_register, payload), nil
}

func (p *Payload) HandleTaskRegister(ctx context.Context, t *asynq.Task) error {
	ti := time.Now()
	var task TaskRegister
	err := sonic.Unmarshal(t.Payload(), &task)
	if err != nil {
		log.Err(err).
			Str("type", "Register").
			Str("payload", string(t.Payload())).
			Msg("Unmarshal error")
		return nil
	}

	log.Info().
		Str("uuid", task.UUID).
		Str("type", "Register").
		Uint("ValidatorID", task.ValidatorID).
		Msg("task start")

	cluster, err := p.getCluster(ctx, p.defaultOperatorStr)
	if err != nil {
		log.Err(err).Msg("scanCluster error")
		p.store.CloseValidatorTask(task.UUID, errors.Wrap(err, "scanCluster error").Error())
		p.store.TurnBackValidatorState(task.ValidatorID, 0)
		return nil
	}

	if !cluster.Active {
		log.Warn().Any("cluster", cluster).Msg("doReactive skip: the cluster is not active")
		p.store.CloseValidatorTask(task.UUID, "the cluster is not active")
		p.store.TurnBackValidatorState(task.ValidatorID, 0)
		return nil
	}

	var keystore *models.ValidatorInfo
	if keystores, err := p.store.GetValidatorStateByID(task.ValidatorID); err != nil {
		log.Err(err).Str("type", "Register").Msg("GetValidatorStateByID error")
		p.store.CloseValidatorTask(task.UUID, errors.Wrap(err, "GetValidatorStateByID error").Error())
		p.store.TurnBackValidatorState(task.ValidatorID, 0)
		return nil
	} else {
		if len(keystores) != 1 {
			log.Error().
				Uint("ValidatorID", task.ValidatorID).
				Str("type", "Register").
				Msg("Validator info not found")
			p.store.CloseValidatorTask(task.UUID, errors.Wrap(err, "Validator info not found").Error())
			p.store.TurnBackValidatorState(task.ValidatorID, 0)
			return nil
		}
		keystore = keystores[0]
	}

	hashBytes := crypto.Keccak256([]byte(keystore.Keystore))
	keystoreSecretKey, err := p.store.GetEncryptedPassword(common.Bytes2Hex(hashBytes))
	if err != nil {
		log.Err(err).
			Uint("ValidatorID", task.ValidatorID).
			Str("type", "Register").
			Msg("GetEncryptedPassword")
		p.store.CloseValidatorTask(task.UUID, errors.Wrap(err, "GetEncryptedPassword not found").Error())
		p.store.TurnBackValidatorState(task.ValidatorID, 0)
		return nil
	}

	password, err := AesDecrypt(keystoreSecretKey, p.keystoreSecretKey)
	if err != nil {
		log.Err(err).
			Uint("ValidatorID", task.ValidatorID).
			Str("type", "Register").
			Msg("AesDecrypt error")
		p.store.CloseValidatorTask(task.UUID, errors.Wrap(err, "AesDecrypt error").Error())
		p.store.TurnBackValidatorState(task.ValidatorID, 0)
		return nil
	}

	nonce := p.store.GetSSVRegisterNonce(ctx)

	cmd := exec.Command(p.makeShares,
		"-op", p.sharesOperatorInfo,
		"-k", keystore.Keystore,
		"-p", password,
		"-o", p.kms.PubKey.String(),
		"-n", nonce)
	shares, err := cmd.CombinedOutput()
	if err != nil {
		log.Err(err).
			Uint("ValidatorID", task.ValidatorID).
			Str("type", "Register").
			Msg("makeShares error")
		p.store.CloseValidatorTask(task.UUID, errors.Wrap(err, "makeShares error").Error())
		p.store.TurnBackValidatorState(task.ValidatorID, 0)
		return nil
	}

	log.Debug().
		Str("publicKey", keystore.PublicKey).
		Any("operatorIds", p.defaultOperatorIds).
		Any("sharesData", string(shares)).
		Any("amount", p.amountTokenSSV).
		Msg("test")

	//p.ssvContract.RegisterValidator(nil, common.FromHex(keystore.PublicKey), p.defaultOperatorIds, common.FromHex(string(shares)), p.amountTokenSSV, *cluster)
	input, _ := p.ssvAbi.Pack("registerValidator", common.FromHex(keystore.PublicKey), p.defaultOperatorIds, common.FromHex(string(shares)), p.amountTokenSSV, *cluster)
	tx, err := p.kms.SignAndSendTransactionByKMSCtx(ctx, p.ssvAddr, input)
	if err != nil {
		if strings.Contains(err.Error(), "ERC20: insufficient allowance") {
			log.Warn().Interface("errmsg", err).Msg("Insufficient SSV approve value learned by error")
			p.doRefreshSSVApprove()
			return err
		}
		p.store.DelSSVClusterSnapshot(ctx, p.defaultOperatorStr)
		p.mail.SendRegisterValidatorErrorMail(keystore.PublicKey, p.defaultOperatorStr, err)
		p.store.CloseValidatorTask(task.UUID, errors.Wrap(err, "Transaction error").Error())
		p.store.TurnBackValidatorState(task.ValidatorID, 0)
		log.Err(err).Msg("doRegisterTask SignAndSendTransactionByKMSCtx error")
		return nil
	}

	p.updateClusterSnapshot(ctx, tx)

	p.store.IncrSSVRegisterNonce(ctx)
	p.store.FinishValidatorTask(task.UUID, p.defaultOperatorStr, tx.TxHash.Hex())
	p.store.StakingValidatorState(task.ValidatorID, p.defaultOperatorStr)

	log.Info().
		Str("uuid", task.UUID).
		Str("txhash", tx.TxHash.Hex()).
		Str("type", "Register").
		Str("nonce", nonce).
		Str("timeuse", time.Now().Sub(ti).String()).
		Msg("task success")

	return nil
}
