package services

import (
	"github.com/fajarardiyanto/chat-application/config"
	"github.com/fajarardiyanto/chat-application/internal/model"
	"github.com/fajarardiyanto/chat-application/internal/repo"
	"github.com/google/uuid"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type ChatService struct{}

func NewChatService() repo.ChatRepository {
	return &ChatService{}
}

func (s *ChatService) CreateChat(req model.Chat) (*model.Chat, error) {
	req.ID = uuid.NewString()
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()

	if err := config.GetDB().Orm().Debug().Model(&model.Chat{}).Create(&req).Error; err != nil {
		config.GetLogger().Error(err)
		return nil, err
	}

	return &req, nil
}

func (s *ChatService) GetChat(from, to string) ([]model.Chat, error) {
	var res []model.Chat
	if err := config.GetDB().Orm().Debug().Model(&model.Chat{}).Where("from_user = ? AND to_user = ? OR from_user = ? AND to_user = ?", from, to, to, from).Order("created_at desc").Find(&res).Error; err != nil {
		config.GetLogger().Error(err)
		return nil, err
	}

	return res, nil
}

func (s *ChatService) SaveFileChat(r *http.Request) (*model.FileModel, error) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		config.GetLogger().Error(err)
		return nil, err
	}

	f, h, err := r.FormFile("file")
	if err != nil {
		config.GetLogger().Error(err)
		return nil, err
	}
	defer f.Close()

	path := "external/" + filepath.Join(".", "files")
	fullPath := path + "/" + h.Filename

	_ = os.MkdirAll(path, os.ModePerm)
	file, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		config.GetLogger().Error(err)
		return nil, err
	}
	defer file.Close()

	if _, err = io.Copy(file, f); err != nil {
		config.GetLogger().Error(err)
		return nil, err
	}

	res := model.FileModel{
		FileName:  fullPath,
		Extension: filepath.Ext(h.Filename),
	}

	return &res, nil
}
