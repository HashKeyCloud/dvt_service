package process

import (
	"context"
	"math/big"
	"sort"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"

	"DVT_Service/abi/ssv"
	"DVT_Service/abi/token"
	"DVT_Service/awsKms"
	com "DVT_Service/common"
	"DVT_Service/conf"
	"DVT_Service/email"
	"DVT_Service/store"
)

type Payload struct {
	store     *store.Store
	ethClient *ethclient.Client
	kms       *awsKms.Kms
	mail      *email.Email

	makeShares         string
	keystoreSecretKey  string
	defaultOperatorStr string
	sharesOperatorInfo string
	defaultOperatorIds []uint64
	defaultOperators   []*conf.Operator

	ssvContract     *ssv.Contract
	ssvAddr         common.Address
	ssvAbi          *abi.ABI
	eventLimitBlock uint64

	ssvTokenContract *token.Contract
	ssvTokenAddr     common.Address
	ssvTokenAbi      *abi.ABI
	amountTokenSSV   *big.Int
}

func InitPayLoadProcess(m *com.Middleware, cfg *conf.Config) *Payload {
	ssvAddr := common.HexToAddress(cfg.SSV.SsvContractAddr)
	ssvContract, _ := ssv.NewContract(ssvAddr, m.EthClient)
	ssvAbi, _ := abi.JSON(strings.NewReader(ssv.ContractMetaData.ABI))

	ssvTokenAddr := common.HexToAddress(cfg.SSV.SsvTokenContractAddr)
	ssvTokenContract, _ := token.NewContract(ssvTokenAddr, m.EthClient)
	ssvTokenAbi, _ := abi.JSON(strings.NewReader(token.ContractMetaData.ABI))

	AmountTokenSSV, flag := new(big.Int).SetString(cfg.SSV.AmountTokenSSV, 10)
	if !flag {
		panic("startServer - config ssv.AmountTokenSSV invalid")
	}

	ApproveTokenSSV, flag := new(big.Int).SetString(cfg.SSV.ApproveCheckTokenSSV, 10)
	if !flag {
		panic("startServer - config ssv.ApproveCheckTokenSSV invalid")
	}

	if AmountTokenSSV.Cmp(ApproveTokenSSV) == 1 {
		panic("startServer - config ApproveCheckTokenSSV is lower than AmountTokenSSV")
	}

	eventLimitBlock := cfg.SSV.ContractCreationBlock

	s := &Payload{
		store:     m.Store,
		kms:       m.KMS,
		ethClient: m.EthClient,
		mail:      m.Email,

		keystoreSecretKey: cfg.Api.KeystoreSecretKey,
		makeShares:        cfg.SSV.MakeShares,

		ssvContract:     ssvContract,
		ssvAddr:         ssvAddr,
		ssvAbi:          &ssvAbi,
		eventLimitBlock: eventLimitBlock,

		ssvTokenContract: ssvTokenContract,
		ssvTokenAddr:     ssvTokenAddr,
		ssvTokenAbi:      &ssvTokenAbi,
		amountTokenSSV:   AmountTokenSSV,
	}

	err := s.tokenSSVApproveCheck(ApproveTokenSSV)
	if err != nil {
		panic(err)
	}

	s.sharesOperatorInfo, _ = sonic.MarshalString(&m.Operators)
	operatorId := make([]uint64, 0, len(m.Operators))
	for _, op := range m.Operators {
		operatorId = append(operatorId, op.ID)
	}
	sort.Slice(operatorId, func(i, j int) bool {
		return operatorId[i] < operatorId[j]
	})

	ctx := context.Background()
	operatorIdStr, _ := sonic.MarshalString(&operatorId)

	s.defaultOperatorStr = operatorIdStr
	s.defaultOperatorIds = operatorId
	s.defaultOperators = m.Operators

	if !s.store.CheckSSVClusterSnapshot(ctx, operatorIdStr) {
		cluster, err := s.scanCluster(operatorIdStr)
		if err != nil {
			panic(err)
		}
		clusterStr, _ := sonic.MarshalString(&cluster)
		s.store.SetSSVClusterSnapshot(ctx, operatorIdStr, clusterStr)
		log.Info().Str("cluster", clusterStr).Msg("get cluster by scanner")
	} else {
		clusterStr := s.store.GetSSVClusterSnapshot(ctx, operatorIdStr)
		log.Info().Str("cluster", clusterStr).Msg("get cluster by redis")
	}

	if !s.store.CheckSSVRegisterNonce(ctx) {
		scanner, err := s.nonceScanner()
		if err != nil {
			panic(err)
		}
		s.store.SetSSVRegisterNonce(ctx, scanner)
		log.Info().Int("nonce", scanner).Msg("get register nonce by scanner")
	} else {
		nonce := s.store.GetSSVRegisterNonce(ctx)
		log.Info().Str("nonce", nonce).Msg("get register nonce by redis")
	}

	log.Info().Msg("PayLoad Process init success")

	return s
}
