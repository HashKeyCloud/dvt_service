package process

import (
	"context"
	"math/big"
	"sort"

	"github.com/bytedance/sonic"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"

	"DVT_Service/abi/ssv"
)

const (
	DAY   = 7200
	WEAK  = DAY * 7
	MONTH = DAY * 30
)

const (
	clusterDepositedTopic   = "ClusterDeposited"
	clusterWithdrawnTopic   = "ClusterWithdrawn"
	validatorRemovedTopic   = "ValidatorRemoved"
	validatorAddTopic       = "ValidatorAdded"
	clusterLiquidatedTopic  = "ClusterLiquidated"
	clusterReactivatedTopic = "ClusterReactivated"
)

func (p *Payload) getCluster(ctx context.Context, operatorIdsMarshal string) (*ssv.ISSVNetworkCoreCluster, error) {
	if p.store.CheckSSVClusterSnapshot(ctx, operatorIdsMarshal) {
		clusterStr := p.store.GetSSVClusterSnapshot(ctx, operatorIdsMarshal)
		var cluster ssv.ISSVNetworkCoreCluster
		sonic.UnmarshalString(clusterStr, &cluster)
		return &cluster, nil
	} else {
		return p.scanCluster(operatorIdsMarshal)
	}
}

// scanCluster Traverse events on the chain to obtain the latest Cluster information
// If there is an error in the process of obtaining the log, the priority is to narrow the scope of the query,
// only adjust it 4 times, and if it still fails, a log acquisition error will be thrown.
// If no event information is found until the height of the contract creation,
// the default value will be returned according to the SSV requirements.
func (p *Payload) scanCluster(operatorIdsMarshal string) (*ssv.ISSVNetworkCoreCluster, error) {
	number, err := p.ethClient.BlockNumber(context.Background())
	if err != nil {
		return nil, err
	}
	log.Info().Msgf("scanCluster %v", number)
	subBlock := big.NewInt(MONTH)
	endBlock := new(big.Int).SetUint64(number)
	beginBlock := new(big.Int).Sub(endBlock, subBlock)

	query := ethereum.FilterQuery{
		Addresses: []common.Address{
			p.ssvAddr,
		},
		Topics: [][]common.Hash{
			{
				p.ssvAbi.Events[clusterDepositedTopic].ID,
				p.ssvAbi.Events[clusterWithdrawnTopic].ID,
				p.ssvAbi.Events[validatorRemovedTopic].ID,
				p.ssvAbi.Events[validatorAddTopic].ID,
				p.ssvAbi.Events[clusterLiquidatedTopic].ID,
				p.ssvAbi.Events[clusterReactivatedTopic].ID,
			},
			{
				p.kms.PubKey.Hash(),
			},
		},
	}

	for endBlock.Uint64() > p.eventLimitBlock {
		query.FromBlock = beginBlock
		query.ToBlock = endBlock

		logs, err := p.ethClient.FilterLogs(context.Background(), query)
		if err != nil {
			switch subBlock.Int64() {
			case MONTH:
				subBlock.SetInt64(WEAK)
			case WEAK:
				subBlock.SetInt64(DAY)
			default:
				return nil, err
			}
			beginBlock = new(big.Int).Sub(endBlock, subBlock)
			continue
		}

		sort.Slice(logs, func(i, j int) bool {
			return logs[i].BlockNumber > logs[j].BlockNumber
		})

		for _, l := range logs {
			var OperatorIds []uint64
			var cluster ssv.ISSVNetworkCoreCluster

			switch l.Topics[0] {
			case p.ssvAbi.Events[clusterDepositedTopic].ID:
				deposited, _ := p.ssvContract.ParseClusterDeposited(l)
				cluster = deposited.Cluster
				OperatorIds = deposited.OperatorIds
			case p.ssvAbi.Events[clusterWithdrawnTopic].ID:
				withdrawn, _ := p.ssvContract.ParseClusterWithdrawn(l)
				cluster = withdrawn.Cluster
				OperatorIds = withdrawn.OperatorIds
			case p.ssvAbi.Events[validatorRemovedTopic].ID:
				removed, _ := p.ssvContract.ParseValidatorRemoved(l)
				cluster = removed.Cluster
				OperatorIds = removed.OperatorIds
			case p.ssvAbi.Events[validatorAddTopic].ID:
				added, _ := p.ssvContract.ParseValidatorAdded(l)
				cluster = added.Cluster
				OperatorIds = added.OperatorIds
			case p.ssvAbi.Events[clusterLiquidatedTopic].ID:
				liquidated, _ := p.ssvContract.ParseClusterLiquidated(l)
				cluster = liquidated.Cluster
				OperatorIds = liquidated.OperatorIds
			case p.ssvAbi.Events[clusterReactivatedTopic].ID:
				reactivated, _ := p.ssvContract.ParseClusterReactivated(l)
				cluster = reactivated.Cluster
				OperatorIds = reactivated.OperatorIds
			}

			OperatorIdsMarshal, _ := sonic.MarshalString(OperatorIds)
			if OperatorIdsMarshal == operatorIdsMarshal {
				log.Info().Interface("cluster", cluster).Msg(l.TxHash.Hex())
				return &cluster, nil
			}
		}

		endBlock = new(big.Int).Sub(beginBlock, big.NewInt(1))
		beginBlock = new(big.Int).Sub(endBlock, subBlock)
	}

	return &ssv.ISSVNetworkCoreCluster{
		ValidatorCount:  0,
		NetworkFeeIndex: 0,
		Index:           0,
		Balance:         big.NewInt(0),
		Active:          true,
	}, nil
}

func (p *Payload) updateClusterSnapshot(ctx context.Context, receipt *types.Receipt) {
	var OperatorIds []uint64
	var cluster ssv.ISSVNetworkCoreCluster
	for _, l := range receipt.Logs {
		switch l.Topics[0] {
		case p.ssvAbi.Events[clusterDepositedTopic].ID:
			deposited, _ := p.ssvContract.ParseClusterDeposited(*l)
			cluster = deposited.Cluster
			OperatorIds = deposited.OperatorIds
		case p.ssvAbi.Events[clusterWithdrawnTopic].ID:
			withdrawn, _ := p.ssvContract.ParseClusterWithdrawn(*l)
			cluster = withdrawn.Cluster
			OperatorIds = withdrawn.OperatorIds
		case p.ssvAbi.Events[validatorRemovedTopic].ID:
			removed, _ := p.ssvContract.ParseValidatorRemoved(*l)
			cluster = removed.Cluster
			OperatorIds = removed.OperatorIds
		case p.ssvAbi.Events[validatorAddTopic].ID:
			added, _ := p.ssvContract.ParseValidatorAdded(*l)
			cluster = added.Cluster
			OperatorIds = added.OperatorIds
		case p.ssvAbi.Events[clusterLiquidatedTopic].ID:
			liquidated, _ := p.ssvContract.ParseClusterLiquidated(*l)
			cluster = liquidated.Cluster
			OperatorIds = liquidated.OperatorIds
		case p.ssvAbi.Events[clusterReactivatedTopic].ID:
			reactivated, _ := p.ssvContract.ParseClusterReactivated(*l)
			cluster = reactivated.Cluster
			OperatorIds = reactivated.OperatorIds
		default:
			continue
		}
		break
	}

	key, _ := sonic.MarshalString(&OperatorIds)
	clusterStr, _ := sonic.MarshalString(&cluster)
	log.Info().Str("key", key).Str("cluster", clusterStr).Msg("get by tx receipt")

	p.store.SetSSVClusterSnapshot(ctx, key, clusterStr)
}

func (p *Payload) nonceScanner() (int, error) {
	number, err := p.ethClient.BlockNumber(context.Background())
	if err != nil {
		return 0, err
	}
	log.Info().Msgf("nonce Scanner %v", number)

	query := ethereum.FilterQuery{
		FromBlock: new(big.Int).SetUint64(p.eventLimitBlock),
		ToBlock:   new(big.Int).SetUint64(number),
		Addresses: []common.Address{
			p.ssvAddr,
		},
		Topics: [][]common.Hash{
			{
				p.ssvAbi.Events[validatorAddTopic].ID,
			},
			{
				p.kms.PubKey.Hash(),
			},
		},
	}

	logs, err := p.ethClient.FilterLogs(context.Background(), query)
	if err != nil {
		return 0, err
	}

	return len(logs), nil
}
