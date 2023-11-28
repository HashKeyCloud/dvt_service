package cliCmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli"

	"DVT_Service/service"
)

func StartProducerServer(ctx *cli.Context) {
	cfg, middleware := configInit(ctx)

	api := service.InitAPIService(cfg.Api, middleware)

	go api.Run()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(interrupt)
	<-interrupt
	log.Info().Msg("Got interrupt, shutting down...")
}
