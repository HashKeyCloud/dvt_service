package service

import (
	"context"
	"fmt"

	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
)

// clusterNonceGet godoc
//
//	@Summary		cluster owner register nonce
//	@Description	get cluster owner register nonce
//	@Tags			nonce
//	@Accept			json
//	@Produce		json
//	@Success		200	 {object} resultResponse "cluster register nonce"
//	@Router			/ssv/cluster/nonce [get]
func (s *APIService) clusterNonceGet(c *gin.Context) {
	nonce := s.store.GetSSVRegisterNonce(context.Background())

	c.JSON(200, resultResponse{
		Code: 200,
		Msg:  fmt.Sprintf("currect cluster register nonce: %s", nonce),
	})
}

// clusterNoncePut godoc
//
//	@Summary		cluster owner register nonce
//	@Description	set cluster owner register nonce
//	@Tags			nonce
//	@Accept			json
//	@Produce		json
//	@Param			body  body	clusterNoncePutRequest	true	"nonce"
//	@Success		200	 {object} resultResponse "massage"
//	@Router			/ssv/cluster/nonce [put]
func (s *APIService) clusterNoncePut(c *gin.Context) {
	var body clusterNoncePutRequest
	if data, err := c.GetRawData(); err != nil || sonic.Unmarshal(data, &body) != nil {
		c.JSON(400, resultResponse{
			Code: 400,
			Msg:  "invalid body",
		})
		return
	}

	s.store.SetSSVRegisterNonce(context.Background(), body.Nonce)

	c.JSON(200, resultResponse{
		Code: 200,
		Msg:  fmt.Sprintf("update cluster register nonce success"),
	})
}
