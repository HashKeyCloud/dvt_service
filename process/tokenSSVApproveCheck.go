package process

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common/math"
	"github.com/rs/zerolog/log"
)

// tokenSSVApproveCheck Confirm the Allowance of the token SSV contract
func (p *Payload) tokenSSVApproveCheck(approveTokenSSV *big.Int) error {
	allowance, _ := p.ssvTokenContract.Allowance(nil, p.kms.PubKey, p.ssvAddr)
	if approveTokenSSV.Cmp(allowance) == 1 {
		log.Warn().Msg("Found that the token SSV allowance is less than the predetermined amount")
		return p.doRefreshSSVApprove()
	}
	return nil
}

// doRefreshSSVApprove processing approve tasks asynchronously
func (p *Payload) doRefreshSSVApprove() error {
	// p.ssvTokenContract.Approve(nil, p.ssvAddr, p.approveTokenSSV)
	input, _ := p.ssvTokenAbi.Pack("approve", p.ssvAddr, math.MaxBig256)
	ctx := context.Background()
	tx, err := p.kms.SignAndSendTransactionByKMSCtx(ctx, p.ssvTokenAddr, input)
	if err != nil {
		log.Err(err).Msg("doRefreshSSVApprove error")
		return err
	}

	log.Info().Msgf("Transaction status check tx: %s is success", tx.TxHash.Hex())
	return nil
}
