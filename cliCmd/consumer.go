package cliCmd

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli"

	"DVT_Service/process"
)

var Consumer = cli.Command{
	Name:        "Consumer",
	ShortName:   "c",
	Usage:       "dvt-service Consumer Server",
	Description: "dvt-service Consumer Server",
	Action:      startConsumerServer,
}

func startConsumerServer(ctx *cli.Context) {
	cfg, middleware := configInit(ctx)

	srv := asynq.NewServer(
		middleware.QueueRedis,
		asynq.Config{
			Concurrency: 1,
			IsFailure: func(err error) bool {
				log.Err(err).Msg("asynq error")
				return false
			},
			Queues: map[string]int{
				"recipient": 7,
				"reactive":  6,
				"deposit":   5,
				"register":  4,
				"remove":    3,
				"withdraw":  2,
				"default":   1,
			},
			StrictPriority: true,
		},
	)

	p := process.InitPayLoadProcess(middleware, cfg)
	go RunHealthy(cfg.Api.Port)

	mux := asynq.NewServeMux()
	mux.HandleFunc(process.Task_SetRecipient, p.HandleTaskSetRecipient)
	mux.HandleFunc(process.Task_Reactive, p.HandleTaskReactive)
	mux.HandleFunc(process.Task_Deposit, p.HandleTaskDeposit)
	mux.HandleFunc(process.Task_register, p.HandleTaskRegister)
	mux.HandleFunc(process.Task_Remove, p.HandleTaskRemove)
	mux.HandleFunc(process.Task_Withdraw, p.HandleTaskWithdraw)
	if err := srv.Run(mux); err != nil {
		log.Fatal().Msgf("could not run server: %v", err)
	}
}

func RunHealthy(port uint) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/healthy", func(c *gin.Context) {
		c.JSON(200, map[string]string{
			"Msg": "healthy",
		})
	})

	log.Info().Msgf("APIHealthy init success, Port: %v", port)
	r.Run(fmt.Sprintf(":%v", port))
}
