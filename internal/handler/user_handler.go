package handler

import (
	"encoding/json"
	"fmt"
	"github.com/fajarardiyanto/chat-application/config"
	"github.com/fajarardiyanto/chat-application/internal/model"
	"github.com/fajarardiyanto/chat-application/internal/model/constant"
	"github.com/fajarardiyanto/chat-application/internal/repo"
	"github.com/fajarardiyanto/chat-application/pkg/auth"
	"github.com/fajarardiyanto/chat-application/pkg/utils"
	"net/http"
	"sync"
)

type UserHandler struct {
	sync.Mutex
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
	u := &model.UserReqModel{}
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
	token, err := auth.ExtractTokenID(r)
	if err != nil {
		config.GetLogger().Error(err)
		model.MessageError(w, http.StatusInternalServerError, err.Error())
		return
	}

	res, err := s.repo.GetUser()
	if err != nil {
		model.MessageError(w, http.StatusNotFound, "no contacts found!")
		return
	}

	var wg sync.WaitGroup
	data := make([]model.UserModel, 0)

	for _, v := range res {
		wg.Add(1)

		go func(wg *sync.WaitGroup, v model.UserModel) {
			defer wg.Done()

			if v.ID != token.ID {
				s.Lock()
				data = append(data, v)
				s.Unlock()
			}
		}(&wg, v)
	}
	wg.Wait()

	model.MessageSuccess(w, http.StatusOK, res)
}

func (s *UserHandler) UpdateStatusHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.ExtractTokenID(r)
	if err != nil {
		config.GetLogger().Error(err)
		model.MessageError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var status bool = false

	userLife := s.repo.CheckUserLife(token.ID)
	if !userLife {
		s.Lock()
		status = true
		s.Unlock()
	}

	if err = s.repo.UpdateStatus(token.ID, status); err != nil {
		model.MessageError(w, http.StatusInternalServerError, err.Error())
		return
	}

	msg := fmt.Sprintf("Successfully update user to %v\n", constant.StatusActivity[status])

	model.MessageSuccessText(w, http.StatusOK, msg)
}
