package service

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"

	"DVT_Service/models"
	"DVT_Service/process"
)

// setFeeRecipientAddress godoc
//
//	@Summary		setFeeRecipientAddress
//	@Description	set FeeRecipient Address on ssv network
//	@Tags			FeeRecipient
//	@Accept			json
//	@Produce		json
//	@Param			body  body	setFeeRecipientAddressBody	true	"New Fee Recipient Address"
//	@Success		200	 {object} resultResponse "FeeRecipientTask Id:(number)"
//	@Failure		400	 {object} resultResponse "params error"
//	@Failure		500	 {object} resultResponse "server error"
//	@Router			/ssv/setFeeRecipientAddress [post]
func (s *APIService) setFeeRecipientAddress(c *gin.Context) {
	var p setFeeRecipientAddressBody
	err := c.BindJSON(&p)
	if err != nil {
		log.Err(err).Str("API", "setFeeRecipientAddress").Msg("invalid body")
		c.JSON(400, resultResponse{
			Code: 400,
			Msg:  "invalid body",
		})
		return
	}

	if len(p.FeeRecipient) == 0 || !common.IsHexAddress(p.FeeRecipient) {
		log.Err(err).Str("API", "setFeeRecipientAddress").Str("Address", p.FeeRecipient).Msg("invalid address")
		c.JSON(400, resultResponse{
			Code: 400,
			Msg:  "invalid address",
		})
		return
	}

	uid := uuid.New().String()
	task := &models.FeeRecipientTask{
		UUID:         uid,
		FeeRecipient: p.FeeRecipient,
	}

	if err := s.store.CreateFeeRecipientTask(task); err != nil {
		log.Err(err).Str("API", "SetFeeRecipient").Msg("CreateFeeRecipientTask error")
		c.JSON(500, resultResponse{
			Code: 500,
			Msg:  "Create SetFeeRecipient Task error",
		})
		return
	}

	depositTask, _ := process.NewSetRecipientTask(uid, p.FeeRecipient)
	if _, err := s.producer.Enqueue(depositTask, asynq.Queue("recipient")); err != nil {
		log.Err(err).Str("API", "SetFeeRecipient").Msg("Enqueue SetFeeRecipient error")
		c.JSON(500, resultResponse{
			Code: 500,
			Msg:  "Enqueue SetFeeRecipient Task error",
		})
		return
	}

	c.JSON(200, resultResponse{
		Code: 200,
		Msg:  fmt.Sprintf("FeeRecipientTask Id: %v", task.UUID),
	})
}
