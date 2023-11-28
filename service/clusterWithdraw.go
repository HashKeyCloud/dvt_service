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

// clusterWithdraw godoc
//
//	@Summary		cluster Withdraw ssv
//	@Description	Withdraw cluster ssv
//	@Tags			cluster
//	@Accept			json
//	@Produce		json
//	@Param			body  body	clusterAmountRequest	true	"amount"
//	@Success		200	 {object} resultResponse "Task uuid"
//	@Failure		400	 {object} resultResponse "params error"
//	@Failure		500	 {object} resultResponse "server error"
//	@Router			/ssv/cluster/withdraw [post]
func (s *APIService) clusterWithdraw(c *gin.Context) {
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
		Type:      3,
		Operators: "",
		Amount:    amount.String(),
	}

	if err := s.store.CreateClusterAmountTask(task); err != nil {
		log.Err(err).Str("API", "clusterWithdraw").Msg("CreateClusterAmountTask error")
		c.JSON(500, resultResponse{
			Code: 500,
			Msg:  "Create Withdraw Task error",
		})
		return
	}

	depositTask, _ := process.NewWithdrawTask(uid, amount.String())
	if _, err := s.producer.Enqueue(depositTask, asynq.Queue("withdraw")); err != nil {
		log.Err(err).Str("API", "clusterWithdraw").Msg("Enqueue Withdraw error")
		c.JSON(500, resultResponse{
			Code: 500,
			Msg:  "Enqueue Withdraw Task error",
		})
		return
	}

	c.JSON(200, resultResponse{
		Code: 200,
		Msg:  fmt.Sprintf("ClusterAmountTask Id: %s", task.UUID),
	})
}
