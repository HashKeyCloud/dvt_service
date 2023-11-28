package conf

// Config DVT service running initialization configuration information structure
type Config struct {
	Api    *ConfigApi   `json:"api"`
	DB     *ConfigDB    `json:"db"`
	SSV    *SSV         `json:"ssv"`
	RPCURL string       `json:"rpc_url"`
	KMS    *KMS         `json:"kms"`
	Email  *ConfigEmail `json:"email"`
}

// ConfigApi Used to configure the gin service
type ConfigApi struct {
	Port              uint   `json:"port"`
	CloseSwagger      bool   `json:"close_swagger,omitempty"`
	KeystoreSecretKey string `json:"keystore_secret_key"`
}

type ConfigDB struct {
	Mysql *ConfigMysql `json:"mysql"`
	Redis *ConfigRedis `json:"redis"`
}

// ConfigMysql Used to configure the mysql service
type ConfigMysql struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
	DBName   string `json:"db"`
}

// ConfigRedis Used to configure the redis service
type ConfigRedis struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

// ConfigEmail Used to configure the email service
type ConfigEmail struct {
	Account  string   `json:"account"`
	Password string   `json:"password"`
	Smtp     string   `json:"smtp"`
	To       []string `json:"to"`
	CC       []string `json:"cc"`
}

// SSV Used to configure SSV network related information
type SSV struct {
	SsvContractAddr       string      `json:"ssv_contract_addr"`
	ContractCreationBlock uint64      `json:"contract_creation_block"`
	SsvTokenContractAddr  string      `json:"ssv_token_contract_addr"`
	AmountTokenSSV        string      `json:"amount_token_ssv"`
	ApproveCheckTokenSSV  string      `json:"approve_check_token_ssv"`
	MakeShares            string      `json:"make_shares"`
	Operators             []*Operator `json:"operators"`
}

type Operator struct {
	ID          uint64 `json:"id"`
	OperatorKey string `json:"operatorKey"`
}

// KMS Used to configure the AWS KMS service
type KMS struct {
	AccessKeyId     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	Region          string `json:"region"`
	KeyId           string `json:"key_id"`
}
