package service

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"

	"DVT_Service/models"
	"DVT_Service/process"
)

// clusterReactive godoc
//
//	@Summary		clusterReactive
//	@Description	reactive cluster
//	@Tags			cluster
//	@Accept			json
//	@Produce		json
//	@Success		200	 {object} resultResponse "Task uuid(string)"
//	@Failure		400	 {object} resultResponse "params error"
//	@Failure		500	 {object} resultResponse "server error"
//	@Router			/ssv/cluster/reactive [post]
func (s *APIService) clusterReactive(c *gin.Context) {
	uid := uuid.New().String()
	task := &models.ClusterAmountTask{
		Type:      1,
		UUID:      uid,
		Operators: "",
	}

	if err := s.store.CreateClusterAmountTask(task); err != nil {
		log.Err(err).Str("API", "clusterReactive").Msg("CreateClusterReactiveTask error")
		c.JSON(500, resultResponse{
			Code: 500,
			Msg:  "Create Reactive Task error",
		})
		return
	}

	reactiveTask, _ := process.NewReactiveTask(uid)
	if _, err := s.producer.Enqueue(reactiveTask, asynq.Queue("reactive")); err != nil {
		log.Err(err).Str("API", "clusterReactive").Msg("Enqueue ReactiveTask error")
		c.JSON(500, resultResponse{
			Code: 500,
			Msg:  "Enqueue Reactive Task error",
		})
		return
	}

	c.JSON(200, resultResponse{
		Code: 200,
		Msg:  fmt.Sprintf("ClusterActiveTask Id: %s", uid),
	})
}
