package utils

import (
	"context"
	"github.com/fajarardiyanto/chat-application/config"
	"github.com/fajarardiyanto/chat-application/internal/model"
	"github.com/fajarardiyanto/flt-go-database/interfaces"
)

func OnMsg(key string, chat *model.Chat) {
	go func() {
		config.GetRabbitMQ().Producer(interfaces.RabbitMQOptions{NoWait: true})
	}()

	go func() {
		if err := config.GetRabbitMQ().Push(context.Background(), "", key, chat, nil); err != nil {
			config.GetLogger().Error(err)
		}
	}()
}
