package service

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// validatorState godoc
//
//	@Summary		validatorState
//	@Description	check Validator state
//	@Tags			Validator
//	@Accept			application/json
//	@Produce		application/json
//	@Param			Validator  path	string	true	"Validator publicKey"
//	@Success		200	 {object} resultResponse "success"
//	@Failure		400	 {object} resultResponse "validator error"
//	@Failure		500	 {object} resultResponse "other fail"
//	@Router			/ssv/{Validator}/state/ [get]
func (s *APIService) validatorState(c *gin.Context) {
	Validator := c.Param("Validator")
	if pub := common.FromHex(Validator); len(pub) != 48 {
		log.Error().Str("API", "validatorState").Str("Validator", Validator).Msg("invalid validator")
		c.JSON(400, resultResponse{
			Code: 400,
			Msg:  "invalid validator",
		})
		return
	} else {
		Validator = common.Bytes2Hex(pub)
	}

	state, err := s.store.GetValidatorStateByPublicKey(Validator)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(400, resultResponse{
				Code: 400,
				Msg:  "validator not found",
			})
		} else {
			log.Err(err).Str("API", "validatorState").Str("Validator", Validator).Msg("validatorState mysql error")
			c.JSON(500, resultResponse{
				Code: 500,
				Msg:  "network error",
			})
		}
		return
	}

	c.JSON(200, resultResponse{
		Code: 200,
		Msg:  "success",
		Data: state.Process(),
	})
}
