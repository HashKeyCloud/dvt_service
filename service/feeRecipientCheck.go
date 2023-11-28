package service

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// feeRecipientState godoc
//
//	@Summary		feeRecipientState
//	@Description	check FeeRecipient ClusterValidatorTask State
//	@Tags			FeeRecipient
//	@Accept			json
//	@Produce		json
//	@Param			TaskId  path	uint	true	"fee recipient task id"
//	@Success		200	 {object} resultResponse "success"
//	@Failure		400	 {object} resultResponse "params error"
//	@Failure		500	 {object} resultResponse "server error"
//	@Router			/ssv/feeRecipient/{TaskId}/state [get]
func (s *APIService) feeRecipientState(c *gin.Context) {
	taskId := c.Param("TaskId")

	tasks, err := s.store.SearchFeeRecipientTaskById(taskId)
	if err != nil {
		log.Err(err).Msg("SearchFeeRecipientTaskById error")
		c.JSON(500, resultResponse{
			Code: 500,
			Msg:  "network error",
		})
		return
	}

	if len(tasks) != 1 {
		c.JSON(400, resultResponse{
			Code: 400,
			Msg:  "task not found",
		})
		return
	}

	c.JSON(200, resultResponse{
		Code: 200,
		Msg:  "success",
		Data: tasks[0],
	})
}
