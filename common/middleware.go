package common

import (
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"

	"DVT_Service/awsKms"
	"DVT_Service/conf"
	"DVT_Service/email"
	"DVT_Service/store"
)

// Middleware collection of middleware
type Middleware struct {
	Store      *store.Store
	Email      *email.Email
	KMS        *awsKms.Kms
	QueueRedis *asynq.RedisClientOpt
	EthClient  *ethclient.Client
	Operators  []*conf.Operator
}

// InitMiddleware Initialize all middleware
func InitMiddleware(c *conf.Config) *Middleware {
	var m Middleware

	m.Store = store.NewStore(c.DB)

	dial, err := ethclient.Dial(c.RPCURL)
	if err != nil {
		panic(fmt.Sprintf("startServer - ethclient create failed!,err: %v", err))
	}

	m.KMS = awsKms.InitKMS(c.KMS, dial)
	log.Info().Str("kms", m.KMS.PubKey.Hex()).Msg("show kms address")

	m.Email = email.InitEmail(c.Email)

	m.QueueRedis = &asynq.RedisClientOpt{
		Addr:     c.DB.Redis.Addr,
		Password: c.DB.Redis.Password,
		DB:       c.DB.Redis.DB,
	}

	m.EthClient = dial
	m.Operators = c.SSV.Operators

	return &m
}
