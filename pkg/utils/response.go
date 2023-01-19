package utils

import (
	"encoding/json"
	"github.com/fajarardiyanto/chat-application/config"
	"github.com/fajarardiyanto/chat-application/internal/model"
	"net/http"
)

func MessageSuccess(w http.ResponseWriter, status int, data interface{}) {
	w.WriteHeader(status)
	res := &model.Response{
		Status: true,
		Data:   data,
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		config.GetLogger().Error(err)
		return
	}
}

func MessageSuccessText(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	res := &model.Response{
		Status:  true,
		Message: msg,
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		config.GetLogger().Error(err)
		return
	}
}

func MessageError(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	res := &model.Response{
		Status:  false,
		Message: msg,
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		config.GetLogger().Error(err)
		return
	}
}
