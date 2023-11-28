package service

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"

	"DVT_Service/models"
	"DVT_Service/process"
)

// registerValidator godoc
//
//	@Summary		registerValidator
//	@Description	Provide the publicKey array to register Validator in ssv network
//	@Tags			Validator
//	@Accept			application/json
//	@Produce		application/json
//	@Param			body  body	validatorBody	true	"publicKey array"
//	@Success		200	 {object} resultResponse "tasks create"
//	@Failure		400	 {object} resultResponse "params error"
//	@Failure		500	 {object} resultResponse "server error"
//	@Router			/ssv/registerValidator [post]
func (s *APIService) registerValidator(c *gin.Context) {
	var p validatorBody
	err := c.BindJSON(&p)
	if err != nil {
		log.Err(err).Str("API", "registerValidator").Msg("invalid body")
		c.JSON(400, resultResponse{
			Code: 400,
			Msg:  "invalid body",
		})
		return
	}

	if len(p.PublicKeys) == 0 {
		log.Error().Str("API", "registerValidator").Interface("Body", p).Msg("without publicKey info")
		c.JSON(400, resultResponse{Code: 400, Msg: "without publicKey info"})
		return
	}

	for i, key := range p.PublicKeys {
		if pub := common.FromHex(key); len(pub) != 48 {
			log.Error().Str("API", "registerValidator").Str("PublicKey", key).Msg("invalid validator")
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
		log.Err(err).Str("API", "registerValidator").Interface("Body", p.PublicKeys).Msg("mysql GetValidatorInfoByPublicKey error")
		c.JSON(500, resultResponse{Code: 500, Msg: "mysql error"})
		return
	}

	if len(keys) != len(p.PublicKeys) {
		log.Err(err).Str("API", "registerValidator").Interface("Body", p.PublicKeys).Int("dbLen", len(keys)).Msg("some publicKey info not found")
		c.JSON(409, resultResponse{Code: 409, Msg: "some publicKey info not found"})
		return
	}

	taskTime := time.Now()
	tasks := make([]*models.ClusterValidatorTask, 0, len(keys))
	registerTasks := make([]*asynq.Task, 0, len(keys))
	for _, key := range keys {
		if key.State != 0 {
			c.JSON(409, resultResponse{
				Code: 409,
				Msg:  fmt.Sprintf("pubkey(%s) is registered", key.PublicKey),
			})
			return
		}

		if key.PendingTime > taskTime.UTC().Unix() {
			errMsg := fmt.Sprintf("publicKey(%s) need wait 2 epoch time(about %v) after removeValidator", key.PublicKey, key.PendingTime)
			log.Error().Str("API", "registerValidator").Int64("PendingTime", key.PendingTime).Msg(errMsg)
			c.JSON(409, resultResponse{
				Code: 409,
				Msg:  errMsg,
			})
			return
		}

		task := &models.ClusterValidatorTask{
			UUID:        uuid.New().String(),
			Type:        1,
			ValidatorID: key.ID,
		}
		tasks = append(tasks, task)

		registerTask, _ := process.NewRegisterTask(task.UUID, key.ID)
		registerTasks = append(registerTasks, registerTask)
	}

	err = s.store.CreateValidatorTasks(tasks)
	if err != nil {
		log.Err(err).Str("API", "registerValidator").Interface("param", tasks).Msg("mysql CreateValidatorTasks error")
		c.JSON(500, resultResponse{Code: 500, Msg: "mysql insert task error"})
		return
	}
	err = s.store.PendingValidatorState(keys)
	if err != nil {
		log.Err(err).Str("API", "registerValidator").Interface("param", keys).Msg("mysql PendingValidatorState error")
		c.JSON(500, resultResponse{Code: 500, Msg: "mysql update validator error"})
		return
	}

	for _, task := range registerTasks {
		if _, err := s.producer.Enqueue(task, asynq.Queue("register")); err != nil {
			log.Err(err).Str("API", "registerValidator").Msg("Enqueue registerValidator error")
			c.JSON(500, resultResponse{
				Code: 500,
				Msg:  "Enqueue registerValidator Task error",
			})
			return
		}
	}

	c.JSON(200, resultResponse{Code: 200, Msg: "tasks create"})
}
