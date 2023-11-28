package awsKms

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

var (
	feeCapMultiple   = big.NewInt(10)
	defaultGasTipCap = big.NewInt(1000000000)
)

// SignAndSendTransactionByKMSCtx Sign the transaction details and initiate the transaction by kms
func (k *Kms) SignAndSendTransactionByKMSCtx(ctx context.Context, to common.Address, data []byte) (*types.Receipt, error) {
	nonce, err := k.ethClient.PendingNonceAt(ctx, k.PubKey)
	if err != nil {
		return nil, err
	}
	suggestedGasPrice, err := k.ethClient.SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}
	suggestedGasLimit, err := k.ethClient.EstimateGas(ctx, ethereum.CallMsg{
		From: k.PubKey,
		To:   &to,
		Data: data,
	})
	if err != nil {
		return nil, err
	}

	gasFeeCap := new(big.Int).Mul(suggestedGasPrice, feeCapMultiple)
	var gasTipCap *big.Int
	if defaultGasTipCap.Cmp(gasFeeCap) == 1 {
		gasTipCap, err = k.ethClient.SuggestGasTipCap(ctx)
		if err != nil {
			return nil, err
		}
	} else {
		gasTipCap = defaultGasTipCap
	}

	tx := types.NewTx(&types.DynamicFeeTx{
		Nonce:     nonce,
		GasTipCap: gasTipCap,
		GasFeeCap: gasFeeCap,
		Gas:       suggestedGasLimit + 1e4,
		To:        &to,
		Data:      data,
	})

	transactOpts, err := k.newAwsKmsTransactorCtx(ctx)
	if err != nil {
		return nil, err
	}

	signerTx, err := transactOpts.Signer(transactOpts.From, tx)
	if err != nil {
		return nil, err
	}

	err = k.ethClient.SendTransaction(ctx, signerTx)
	if err != nil {
		return nil, err
	}

	return bind.WaitMined(ctx, k.ethClient, signerTx)
}
