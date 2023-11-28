package awsKms

import (
	"context"
	"encoding/asn1"
	"fmt"
	"math/big"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"

	"DVT_Service/conf"
)

// Kms Construct a kms transaction with and sign it and send it
type Kms struct {
	ethClient *ethclient.Client
	ChainId   *big.Int

	PubKey      common.Address
	pubKeyBytes []byte
	client      *kms.Client
	keyId       string
}

type asn1EcPublicKey struct {
	EcPublicKeyInfo asn1EcPublicKeyInfo
	PublicKey       asn1.BitString
}

type asn1EcPublicKeyInfo struct {
	Algorithm  asn1.ObjectIdentifier
	Parameters asn1.ObjectIdentifier
}

// InitKMS Init Kms struct
func InitKMS(cfg *conf.KMS, dial *ethclient.Client) *Kms {
	chainId, err := dial.ChainID(context.Background())
	if err != nil {
		panic(fmt.Sprintf("startServer - can't get chainId!,err: %v", err))
	}

	awsCfg := aws.Config{
		Credentials: credentials.NewStaticCredentialsProvider(cfg.AccessKeyId, cfg.AccessKeySecret, ""),
		Retryer: func() aws.Retryer {
			return aws.NopRetryer{}
		},
		Region: cfg.Region,
	}
	client := kms.NewFromConfig(awsCfg)

	key, err := client.GetPublicKey(context.Background(), &kms.GetPublicKeyInput{
		KeyId: aws.String(cfg.KeyId),
	})
	if err != nil {
		panic(err)
	}

	var asn1pubk asn1EcPublicKey
	_, err = asn1.Unmarshal(key.PublicKey, &asn1pubk)
	if err != nil {
		panic(errors.Wrapf(err, "can not parse asn1 public key for KeyId=%s", *key.KeyId))
	}

	pubkey, err := crypto.UnmarshalPubkey(asn1pubk.PublicKey.Bytes)
	if err != nil {
		panic(err)
	}

	return &Kms{
		ethClient: dial,
		ChainId:   chainId,

		PubKey:      crypto.PubkeyToAddress(*pubkey),
		pubKeyBytes: secp256k1.S256().Marshal(pubkey.X, pubkey.Y),
		client:      client,
		keyId:       cfg.KeyId,
	}
}
