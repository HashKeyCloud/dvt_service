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

// removeValidator godoc
//
//	@Summary		removeValidator
//	@Description	remove Validator to ssv network
//	@Tags			Validator
//	@Accept			application/json
//	@Produce		application/json
//	@Param			body  body	validatorBody	true	"publicKey array"
//	@Success		200	 {object} resultResponse "tasks create"
//	@Failure		400	 {object} resultResponse "params error"
//	@Failure		500	 {object} resultResponse "server error"
//	@Router			/ssv/removeValidator [post]
func (s *APIService) removeValidator(c *gin.Context) {
	var p validatorBody
	err := c.BindJSON(&p)
	if err != nil {
		log.Err(err).Str("API", "removeValidator").Msg("invalid body")
		c.JSON(400, resultResponse{
			Code: 400,
			Msg:  "invalid body",
		})
		return
	}

	if len(p.PublicKeys) == 0 {
		log.Error().Str("API", "removeValidator").Interface("Body", p).Msg("without publicKey info")
		c.JSON(400, resultResponse{
			Code: 400,
			Msg:  "without publicKey info",
		})
		return
	}

	for i, key := range p.PublicKeys {
		if pub := common.FromHex(key); len(pub) != 48 {
			log.Error().Str("API", "removeValidator").Str("PublicKey", key).Msg("invalid validator")
			c.JSON(400, resultResponse{
				Code: 400,
				Msg:  "invalid validator",
			})
			return
		} else {
			p.PublicKeys[i] = common.Bytes2Hex(pub)
		}
	}

	keys, err := s.store.GetValidatorInfoByPublicKey(p.PublicKeys)
	if err != nil {
		log.Err(err).Str("API", "removeValidator").Interface("params", p.PublicKeys).Msg("mysql GetStakingValidatorInfoByPublicKey error")
		c.JSON(500, resultResponse{Code: 500, Msg: "mysql error"})
		return
	}

	if len(keys) != len(p.PublicKeys) {
		log.Err(err).Str("API", "removeValidator").Interface("Body", p.PublicKeys).Int("dbLen", len(keys)).Msg("some publicKey info not found")
		c.JSON(400, resultResponse{Code: 400, Msg: "some publicKey info not found"})
		return
	}

	tasks := make([]*models.ClusterValidatorTask, 0, len(keys))
	removeTasks := make([]*asynq.Task, 0, len(keys))
	for _, key := range keys {
		if key.State != 2 {
			c.JSON(409, resultResponse{
				Code: 409,
				Msg:  fmt.Sprintf("pubkey(%s) is not registered", key.PublicKey),
			})
			return
		}

		task := &models.ClusterValidatorTask{
			UUID:        uuid.New().String(),
			Type:        2,
			ValidatorID: key.ID,
		}
		tasks = append(tasks, task)

		removeTask, _ := process.NewRemoveTask(task.UUID, task.ValidatorID)
		removeTasks = append(removeTasks, removeTask)
	}

	err = s.store.CreateValidatorTasks(tasks)
	if err != nil {
		log.Err(err).Str("API", "removeValidator").Interface("params", tasks).Msg("mysql CreateValidatorTasks error")
		c.JSON(500, resultResponse{Code: 400, Msg: "mysql insert task error"})
		return
	}

	err = s.store.PendingValidatorState(keys)
	if err != nil {
		log.Err(err).Str("API", "removeValidator").Interface("params", keys).Msg("mysql PendingValidatorState error")
		c.JSON(500, resultResponse{Code: 400, Msg: "mysql update validator error"})
		return
	}

	for _, task := range removeTasks {
		if _, err := s.producer.Enqueue(task, asynq.Queue("remove")); err != nil {
			log.Err(err).Str("API", "removeValidator").Msg("Enqueue removeValidator error")
			c.JSON(500, resultResponse{
				Code: 500,
				Msg:  "Enqueue removeValidator Task error",
			})
			return
		}
	}

	c.JSON(200, resultResponse{Code: 200, Msg: "tasks create"})
}
