package commons

import (
	"github.com/fajarardiyanto/chat-application/config"
	"github.com/fajarardiyanto/chat-application/pkg/auth"
	"net/http"
)

func ExtractTokenRole(r *http.Request) {
	_, err := auth.ExtractTokenID(r)
	if err != nil {
		config.GetLogger().Error(err)
		return
	}
}
