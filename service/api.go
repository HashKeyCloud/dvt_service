package service

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"

	_ "DVT_Service/docs"
	swaggerFiles "github.com/swaggo/files"
	swagger "github.com/swaggo/gin-swagger"

	com "DVT_Service/common"
	"DVT_Service/conf"
	"DVT_Service/store"
)

// APIService go-gin api service, build http service
type APIService struct {
	port              uint
	closeSwagger      bool
	keystoreSecretKey string

	store    *store.Store
	producer *asynq.Client
}

// InitAPIService Init go-gin api service
func InitAPIService(cfg *conf.ConfigApi, middleware *com.Middleware) *APIService {
	return &APIService{
		port:              cfg.Port,
		closeSwagger:      cfg.CloseSwagger,
		keystoreSecretKey: cfg.KeystoreSecretKey,
		store:             middleware.Store,
		producer:          asynq.NewClient(middleware.QueueRedis),
	}
}

// Run Start go-gin api service
func (s *APIService) Run() {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(cors.Default())

	r.GET("/healthy", s.healthy)

	if !s.closeSwagger {
		// swagger plugin
		r.GET("/swagger/*any", swagger.WrapHandler(swaggerFiles.Handler))
	}

	ssvGroup := r.Group("ssv", gin.Logger())
	ssvGroup.POST("/upload", s.upload)

	ssvGroup.POST("/registerValidator", s.registerValidator)
	ssvGroup.POST("/removeValidator", s.removeValidator)
	ssvGroup.GET("/:Validator/state", s.validatorState)

	ssvGroup.POST("/setFeeRecipientAddress", s.setFeeRecipientAddress)
	ssvGroup.GET("/feeRecipient/:TaskId/state", s.feeRecipientState)

	ssvGroup.POST("/cluster/reactive", s.clusterReactive)
	ssvGroup.POST("/cluster/deposit", s.clusterDeposit)
	ssvGroup.POST("/cluster/withdraw", s.clusterWithdraw)

	ssvGroup.GET("/cluster/nonce", s.clusterNonceGet)
	ssvGroup.PUT("/cluster/nonce", s.clusterNoncePut)

	log.Info().Msgf("APIService init success, Port: %v", s.port)
	r.Run(fmt.Sprintf(":%v", s.port))
}
