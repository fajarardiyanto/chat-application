package model

type TokenModel struct {
	UserId    string `json:"userId"`
	AccountId string `json:"accountId"`
	Role      string `json:"role"`
	SessionId string `json:"sessionId"`
}
