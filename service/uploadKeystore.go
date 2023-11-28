package service

import (
	"github.com/bytedance/sonic"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	keystorev4 "github.com/wealdtech/go-eth2-wallet-encryptor-keystorev4"

	"DVT_Service/models"
)

// upload godoc
//
//	@Summary		upload
//	@Description	upload keystore by dvt tools
//	@Tags			upload
//	@Accept			application/json
//	@Produce		application/json
//	@Param			body  body	models.UploadKeystore	true	"keystore info"
//	@Success		200	 {object} resultResponse "success"
//	@Failure		400	 {object} resultResponse "params error"
//	@Failure		500	 {object} resultResponse "server error"
//	@Router			/ssv/upload [post]
func (s *APIService) upload(c *gin.Context) {
	rawData, _ := c.GetRawData()
	var body models.UploadKeystore
	if err := sonic.Unmarshal(rawData, &body); err != nil || len(body.Keys) == 0 || len(body.Password) < 8 {
		c.JSON(400, models.ResultResponse[string]{
			Code: 400,
			Msg:  "error",
			Data: "invalid rawData",
		})
		return
	}

	if _, err := keystorev4.New().Decrypt(body.Keys[0].Crypto, body.Password); err != nil {
		c.JSON(400, models.ResultResponse[string]{
			Code: 400,
			Msg:  "error",
			Data: "invalid password",
		})
		return
	}

	links := make([]string, len(body.Keys))
	keyStrs := make([]string, len(body.Keys))
	pubkeys := make([]string, len(body.Keys))
	encrypts := make([]string, len(body.Keys))

	for i, key := range body.Keys {
		ksv4Bytes, _ := sonic.Marshal(key)
		hash := crypto.Keccak256(ksv4Bytes)
		link := common.Bytes2Hex(hash)
		encrypt, _ := AesEncrypt(body.Password, s.keystoreSecretKey)

		links[i] = link
		keyStrs[i] = string(ksv4Bytes)
		pubkeys[i] = key.Pubkey
		encrypts[i] = encrypt
	}

	if err := s.store.SaveValidator(keyStrs, links, pubkeys, encrypts); err != nil {
		c.JSON(500, models.ResultResponse[string]{
			Code: 500,
			Msg:  "error",
			Data: err.Error(),
		})
		return
	}

	c.JSON(200, models.ResultResponse[string]{
		Code: 200,
		Msg:  "success",
		Data: "upload success",
	})
}
