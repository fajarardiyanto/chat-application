package services

import (
	"github.com/fajarardiyanto/chat-application/config"
	"github.com/fajarardiyanto/chat-application/internal/model"
	"github.com/fajarardiyanto/chat-application/internal/repo"
	"github.com/fajarardiyanto/chat-application/pkg/utils"
	"github.com/google/uuid"
	"time"
)

type UserService struct{}

func NewUserService() repo.UserRepository {
	return &UserService{}
}

func (*UserService) UserExist(username string) error {
	var res model.UserModel
	if err := config.GetDB().Orm().Debug().Model(&model.UserModel{}).Where("username = ?", username).First(&res).Error; err != nil {
		config.GetLogger().Error(err)
		return err
	}

	return nil
}

func (*UserService) Register(req model.UserModel) (*model.UserModel, error) {
	pass, err := utils.HashPassword(req.Password)
	if err != nil {
		config.GetLogger().Error(err)
		return nil, err
	}

	req.ID = uuid.NewString()
	req.Password = pass
	req.Status = false
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()

	if err = config.GetDB().Orm().Debug().Model(&model.UserModel{}).Create(&req).Error; err != nil {
		config.GetLogger().Error(err)
		return nil, err
	}

	return &req, nil
}

func (*UserService) GetUser() ([]model.UserModel, error) {
	var res []model.UserModel
	if err := config.GetDB().Orm().Debug().Model(&model.UserModel{}).Find(&res).Error; err != nil {
		config.GetLogger().Error(err)
		return nil, err
	}

	return res, nil
}

func (*UserService) UpdateStatus(id string, status bool) error {
	if err := config.GetDB().Orm().Debug().Model(model.UserModel{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}).Error; err != nil {
		config.GetLogger().Error(err)
		return err
	}

	return nil
}

func (*UserService) CheckUserLife(id string) bool {
	var res model.UserModel
	if err := config.GetDB().Orm().Debug().Model(&model.UserModel{}).Where("id = ?", id).First(&res).Error; err != nil {
		config.GetLogger().Error(err)
		return false
	}

	return res.Status
}
