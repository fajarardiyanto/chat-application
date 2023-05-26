package model

type AgentTokenModel struct {
	UserId    string `json:"userId"`
	AccountId string `json:"accountId"`
	Role      string `json:"role"`
	SessionId string `json:"sessionId"`
}

type ContactTokenModel struct {
	SourceId string `json:"sourceId"`
	InboxId  string `json:"inboxId"`
}
