package service

// resultResponse General api result return struct
type resultResponse struct {
	Code int         `json:"code"`           // response code
	Msg  string      `json:"msg"`            // error or success msg
	Data interface{} `json:"data,omitempty"` // other data info
}

// setFeeRecipientAddressBody API setFeeRecipientAddress dedicated body receiving struct
type setFeeRecipientAddressBody struct {
	FeeRecipient string `json:"fee_recipient"` // New Fee Recipient Address
}

// validatorBody API registerValidator and removeValidator dedicated body receiving struct
type validatorBody struct {
	PublicKeys []string `json:"public_keys"` // public key array
}

type clusterAmountRequest struct {
	Amount string `json:"amount"`
}

type clusterNoncePutRequest struct {
	Nonce int `json:"nonce"`
}
