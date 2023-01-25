package handler

import (
	"encoding/json"
	"fmt"
	"github.com/fajarardiyanto/chat-application/config"
	"github.com/fajarardiyanto/chat-application/internal/model"
	"github.com/fajarardiyanto/chat-application/internal/repo"
	"github.com/fajarardiyanto/chat-application/pkg/auth"
	"github.com/fajarardiyanto/chat-application/pkg/utils"
	"net/http"
)

type UserHandler struct {
	repo repo.UserRepository
}

func NewUserHandler(repo repo.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

func (s *UserHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	u := model.UserReqModel{}
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		model.MessageError(w, http.StatusBadRequest, "error decoding request object")
		return
	}

	if _, err := s.repo.UserExist(u.Username); err == nil {
		model.MessageError(w, http.StatusFound, "username already exist!")
		return
	}

	req := model.UserModel{
		Username: u.Username,
		Password: u.Password,
		UserType: u.UserType,
	}

	res, err := s.repo.Register(req)
	if err != nil {
		config.GetLogger().Error(err.Error())
		model.MessageError(w, http.StatusInternalServerError, "something went wrong while registering the user. please try again after sometime.")
		return
	}

	model.MessageSuccess(w, http.StatusOK, res)
}

func (s *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	u := &model.UserReq{}
	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		model.MessageError(w, http.StatusBadRequest, "error decoding request object")
		return
	}

	res, err := s.repo.UserExist(u.Username)
	if err != nil {
		config.GetLogger().Error(err)
		model.MessageError(w, http.StatusNotFound, "Invalid username/password")
		return
	}

	if err = utils.VerifyPassword(res.Password, u.Password); err != nil {
		config.GetLogger().Error(err)
		model.MessageError(w, http.StatusBadRequest, "Invalid username/password")
		return
	}

	userToken := model.UserTokenModel{
		ID:       res.ID,
		Username: res.Username,
		UserType: res.UserType,
	}

	token, err := auth.CreateToken(userToken)
	if err != nil {
		config.GetLogger().Error(err)
		model.MessageError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	response := model.UserResponseModel{
		User:  *res,
		Token: token,
	}

	model.MessageSuccess(w, http.StatusOK, response)
}

func (s *UserHandler) ContactListHandler(w http.ResponseWriter, r *http.Request) {
	username := utils.QueryString(r, "username")

	res, err := s.repo.GetUser()
	if err != nil {
		model.MessageError(w, http.StatusNotFound, "no contacts found!")
		return
	}

	data := make([]model.UserModel, 0)
	for _, v := range res {
		if v.Username != username {
			data = append(data, v)
		}
	}

	model.MessageSuccess(w, http.StatusOK, res)
}

func (s *UserHandler) UpdateStatusHandler(w http.ResponseWriter, r *http.Request) {
	id := utils.QueryString(r, "id")

	var status bool = false

	userLife := s.repo.CheckUserLife(id)
	if !userLife {
		status = true
	}

	if err := s.repo.UpdateStatus(id, status); err != nil {
		model.MessageError(w, http.StatusInternalServerError, err.Error())
		return
	}

	msg := fmt.Sprintf("Successfully update user to %v\n", model.StatusActivity[status])

	model.MessageSuccessText(w, http.StatusOK, msg)
}
