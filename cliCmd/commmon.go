package cliCmd

import (
	"fmt"

	"github.com/urfave/cli"

	"DVT_Service/common"
	"DVT_Service/conf"
)

func configInit(ctx *cli.Context) (*conf.Config, *common.Middleware) {
	cfg, err := conf.LoadConfig(ctx)
	if err != nil {
		panic(fmt.Sprintf("startServer - read config failed!,err: %v", err))
	}

	mysqlPassword := ctx.GlobalString(conf.MysqlPassword)
	if len(mysqlPassword) > 0 {
		cfg.DB.Mysql.Password = mysqlPassword
	}

	redisPassword := ctx.GlobalString(conf.RedisPassword)
	if len(redisPassword) > 0 {
		cfg.DB.Redis.Password = redisPassword
	}

	keystoreSecretKey := ctx.GlobalString(conf.KeystoreSecretKey)
	if len(keystoreSecretKey) > 0 {
		cfg.Api.KeystoreSecretKey = keystoreSecretKey
	}

	middleware := common.InitMiddleware(cfg)

	return cfg, middleware
}
