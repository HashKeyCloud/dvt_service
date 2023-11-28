package conf

import (
	"net/url"
	"strconv"

	"github.com/bytedance/sonic"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// loadFromNacos Obtain configuration information through nacos and parse it to the config structure
func loadFromNacos(url, namespace, dataId, group, userName, password string) (*Config, error) {
	connect, err := nacosConnectAndGet(url, namespace, dataId, group, userName, password)
	if err != nil {
		log.Error().Msgf("Connect: failed, err: %s", err)
		return nil, errors.Wrap(err, "failed on connect")
	}

	var cfg Config
	err = sonic.UnmarshalString(connect, &cfg)
	if err != nil {
		log.Error().Msgf("LoadConfig: failed, err: %s", err)
		return nil, errors.Wrap(err, "failed on unmarshal config")
	}

	log.Info().Msg("LoadConfig form nacos success")
	return &cfg, nil
}

// nacosConnectAndGet Connect to nacos to get configuration information
func nacosConnectAndGet(urlstr, namespace, dataId, group, userName, password string) (string, error) {
	nacosURL, err := url.Parse(urlstr)
	if err != nil {
		log.Err(err).Msgf("url.Parse %s", urlstr)
		return "", err
	}
	portString := nacosURL.Port()
	var port uint64 = 80
	if len(portString) == 0 {
		if nacosURL.Scheme == "https" {
			port = 443
		}
	} else {
		port, err = strconv.ParseUint(portString, 10, 64)
		if err != nil {
			log.Err(err).Msgf("ParseUint: failed %s", portString)
			return "", err
		}
	}

	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      nacosURL.Hostname(),
			ContextPath: nacosURL.RequestURI(),
			Scheme:      nacosURL.Scheme,
			Port:        port,
		},
	}

	clientConfig := constant.ClientConfig{
		NamespaceId:         namespace,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "info",
	}

	if userName != "" && password != "" {
		clientConfig.Username = userName
		clientConfig.Password = password
	}

	// Create naming client for service discovery
	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		log.Err(err).Msg("CreateConfigClient: failed")
		return "", err
	}

	// Read
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
	})

	if err != nil {
		log.Err(err).Msg("GetConfig: failed")
		return "", err
	}
	return content, nil
}
