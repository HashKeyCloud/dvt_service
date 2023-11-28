package service

import (
	"fmt"
	"math/big"

	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"

	"DVT_Service/models"
	"DVT_Service/process"
)

// clusterDeposit godoc
//
//	@Summary		cluster Deposit ssv
//	@Description	Deposit cluster ssv
//	@Tags			cluster
//	@Accept			json
//	@Produce		json
//	@Param			body  body	clusterAmountRequest	true	"amount"
//	@Success		200	 {object} resultResponse "Task uuid"
//	@Failure		400	 {object} resultResponse "params error"
//	@Failure		500	 {object} resultResponse "server error"
//	@Router			/ssv/cluster/deposit [post]
func (s *APIService) clusterDeposit(c *gin.Context) {
	var body clusterAmountRequest
	if data, err := c.GetRawData(); err != nil || sonic.Unmarshal(data, &body) != nil {
		c.JSON(400, resultResponse{
			Code: 400,
			Msg:  "invalid body",
		})
		return
	}

	amount, b := new(big.Int).SetString(body.Amount, 10)
	if !b {
		c.JSON(400, resultResponse{
			Code: 400,
			Msg:  "invalid amount",
		})
	}

	uid := uuid.New().String()
	task := &models.ClusterAmountTask{
		UUID:      uid,
		Type:      2,
		Operators: "",
		Amount:    amount.String(),
	}

	if err := s.store.CreateClusterAmountTask(task); err != nil {
		log.Err(err).Str("API", "clusterDeposit").Msg("CreateClusterAmountTask error")
		c.JSON(500, resultResponse{
			Code: 500,
			Msg:  "Create Deposit Task error",
		})
		return
	}

	depositTask, _ := process.NewDepositTask(uid, amount.String())
	if _, err := s.producer.Enqueue(depositTask, asynq.Queue("deposit")); err != nil {
		log.Err(err).Str("API", "clusterDeposit").Msg("Enqueue DepositTask error")
		c.JSON(500, resultResponse{
			Code: 500,
			Msg:  "Enqueue Deposit Task error",
		})
		return
	}

	c.JSON(200, resultResponse{
		Code: 200,
		Msg:  fmt.Sprintf("ClusterAmountTask Id: %s", task.UUID),
	})
}
