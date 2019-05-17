package util

import "encoding/json"

const (
	PublicKey  = "/storage/cert/auth_public_key.pem"
	PrivateKey = "/storage/cert/auth_private_key.pem"
)

// token 数据结构
type TokenData struct {
	Token              json.Token `json:"token"`
	Expire             int64      `json:"expire"`
}
