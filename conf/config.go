package conf

import (
	"os"

	"github.com/bytedance/sonic"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli"
)

const (
	NacosUrl         = "NACOS_URL"
	NacosNamespaceId = "NACOS_NAMESPACE_ID"
	NacosDataId      = "NACOS_DATA_ID"
	NacosGroup       = "NACOS_GROUP"
	NacosUsername    = "NACOS_USERNAME"
	NacosPassword    = "NACOS_PASSWORD"

	KeystoreSecretKey = "KEYSTORE_SECRET_KEY"
	MysqlPassword     = "MYSQL_PASSWORD"
	RedisPassword     = "REDIS_PASSWORD"
)

// LoadConfig Load config info
// Determine whether the variable NACOS_URL has a value,
// and choose whether to obtain it through nacos
// or read it through a configuration file.
func LoadConfig(ctx *cli.Context) (*Config, error) {
	nacosUrl := ctx.GlobalString(NacosUrl)
	if len(nacosUrl) > 0 {
		nacosNamespaceId := ctx.GlobalString(NacosNamespaceId)
		nacosDataId := ctx.GlobalString(NacosDataId)
		nacosGroup := ctx.GlobalString(NacosGroup)
		nacosUsername := ctx.GlobalString(NacosUsername)
		nacosPassword := ctx.GlobalString(NacosPassword)
		return loadFromNacos(nacosUrl, nacosNamespaceId, nacosDataId, nacosGroup, nacosUsername, nacosPassword)
	}

	path := ctx.GlobalString("c")
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = sonic.Unmarshal(file, &cfg)
	if err != nil {
		return nil, err
	}

	log.Info().Msg("LoadConfig form config file success")
	return &cfg, nil
}
