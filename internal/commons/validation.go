package commons

import (
	"github.com/fajarardiyanto/chat-application/config"
	"github.com/fajarardiyanto/chat-application/internal/constant"
	"github.com/fajarardiyanto/chat-application/pkg/auth"
	"net/http"
)

func IsAllowedToSetPassword(r *http.Request) bool {
	return AllowedToPerformAction(r, []int32{int32(constant.SUPERVISOR), int32(constant.ADMIN)})
}

func IsAllowedToGetAgent(r *http.Request) bool {
	return AllowedToPerformAction(r, []int32{int32(constant.SUPERVISOR), int32(constant.AGENT)})
}

func IsAllowedToGetAllAgent(r *http.Request) bool {
	return AllowedToPerformAction(r, []int32{int32(constant.SUPERVISOR), int32(constant.ADMIN)})
}

func IsAllowedToDeleteAgent(r *http.Request) bool {
	return AllowedToPerformAction(r, []int32{int32(constant.SUPERVISOR), int32(constant.ADMIN)})
}

func IsAllowedToChat(r *http.Request) bool {
	return AllowedToPerformAction(r, []int32{int32(constant.AGENT)})
}

func AllowedToPerformAction(r *http.Request, roles []int32) bool {
	token, err := auth.ExtractTokenID(r)
	if err != nil {
		config.GetLogger().Error(err)
		return false
	}

	for _, v := range roles {
		if constant.AgentRole(v) == token.Role {
			return true
		}
	}
	return false
}
