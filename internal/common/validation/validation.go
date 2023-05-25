package validation

import (
	"github.com/fajarardiyanto/chat-application/config"
	"github.com/fajarardiyanto/chat-application/internal/common/constant"
	"github.com/fajarardiyanto/chat-application/pkg/auth"
	"net/http"
)

func IsAllowedToSetPassword(r *http.Request) bool {
	return AllowedToPerformAction(r, []constant.Role{constant.SUPERVISOR, constant.ADMIN})
}

func IsAllowedToGetAgent(r *http.Request) bool {
	return AllowedToPerformAction(r, []constant.Role{constant.SUPERVISOR, constant.AGENT})
}

func IsAllowedToGetAllAgent(r *http.Request) bool {
	return AllowedToPerformAction(r, []constant.Role{constant.SUPERVISOR, constant.ADMIN})
}

func IsAllowedToDeleteAgent(r *http.Request) bool {
	return AllowedToPerformAction(r, []constant.Role{constant.SUPERVISOR, constant.ADMIN})
}

func IsAllowedToChat(r *http.Request) bool {
	return AllowedToPerformAction(r, []constant.Role{constant.AGENT})
}

func AllowedToPerformAction(r *http.Request, roles []constant.Role) bool {
	token, err := auth.ExtractTokenID(r)
	if err != nil {
		config.GetLogger().Error(err)
		return false
	}

	for _, v := range roles {
		if constant.AgentRole(int32(v)) == token.Role {
			return true
		}
	}
	return false
}
