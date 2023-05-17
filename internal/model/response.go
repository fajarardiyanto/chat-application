package model

import (
	"encoding/json"
	"github.com/fajarardiyanto/chat-application/config"
	"net/http"
)

type Exception struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func MessageSuccess(w http.ResponseWriter, status int, data interface{}) {
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		config.GetLogger().Error(err)
		return
	}
}

func MessageError(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	res := &Exception{
		Status:  status,
		Message: msg,
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		config.GetLogger().Error(err)
		return
	}
}
