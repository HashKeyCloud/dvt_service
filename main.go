package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/urfave/cli"

	"DVT_Service/cliCmd"
	"DVT_Service/common"
	"DVT_Service/conf"
)

func setupApp() *cli.App {
	app := cli.NewApp()
	app.Usage = "DVT Service"
	app.Action = cliCmd.StartProducerServer
	app.Version = "1.0.0"
	app.Commands = []cli.Command{
		cliCmd.Consumer,
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config, c",
			EnvVar: "CONFIG_FILE_PATH",
			Usage:  "Service config file `<path>`",
			Value:  "",
		},
		// NACOS
		cli.StringFlag{
			Name:   conf.NacosUrl,
			EnvVar: conf.NacosUrl,
		},
		cli.StringFlag{
			Name:   conf.NacosNamespaceId,
			EnvVar: conf.NacosNamespaceId,
		},
		cli.StringFlag{
			Name:   conf.NacosDataId,
			EnvVar: conf.NacosDataId,
		},
		cli.StringFlag{
			Name:   conf.NacosGroup,
			EnvVar: conf.NacosGroup,
			Value:  "DEFAULT_GROUP",
		},
		cli.StringFlag{
			Name:   conf.NacosUsername,
			EnvVar: conf.NacosUsername,
		},
		cli.StringFlag{
			Name:   conf.NacosPassword,
			EnvVar: conf.NacosPassword,
		},
		// Mysql Password
		cli.StringFlag{
			Name:   conf.MysqlPassword,
			EnvVar: conf.MysqlPassword,
		},
		// Redis Password
		cli.StringFlag{
			Name:   conf.RedisPassword,
			EnvVar: conf.RedisPassword,
		},
		// Keystore Secret Key
		cli.StringFlag{
			Name:   conf.KeystoreSecretKey,
			EnvVar: conf.KeystoreSecretKey,
		},
	}

	app.Before = func(context *cli.Context) error {
		cpuNum := 0
		if runtime.NumCPU() > 4 {
			cpuNum = 4
		} else {
			cpuNum = runtime.NumCPU()
		}
		runtime.GOMAXPROCS(cpuNum)
		return nil
	}
	return app
}

//	@title			DVT Service API
//	@version		1.0.0
//	@description	This is a sample server celler server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	common.InitLogger("logs/dvt_service.log")

	if err := setupApp().Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
