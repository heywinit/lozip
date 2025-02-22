package models

type NonceRequest struct {
	PublicKey string `json:"publicKey" binding:"required"`
}

type VerifyRequest struct {
	PublicKey string `json:"publicKey" binding:"required"`
	Signature string `json:"signature" binding:"required"`
}

type LogoutRequest struct {}

