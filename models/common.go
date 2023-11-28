package models

// ResultResponse General api result return struct
type ResultResponse[T any] struct {
	Code int    `json:"code"`           // response code
	Msg  string `json:"msg"`            // error or success msg
	Data T      `json:"data,omitempty"` // data info
}

type UploadKeystore struct {
	Keys     []*KeystoreV4 `json:"keys"`
	Password string        `json:"password"`
}

type KeystoreV4 struct {
	Crypto  map[string]interface{} `json:"crypto"`
	Uuid    string                 `json:"uuid"`
	Pubkey  string                 `json:"pubkey"`
	Version int                    `json:"version"`
	Name    string                 `json:"name"`
	Path    string                 `json:"path"`
}
