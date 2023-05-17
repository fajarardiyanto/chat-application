package controller

import (
	"encoding/json"
	"github.com/fajarardiyanto/chat-application/config"
	"github.com/fajarardiyanto/chat-application/internal/constant"
	"github.com/fajarardiyanto/chat-application/internal/controller/dto/request"
	"github.com/fajarardiyanto/chat-application/internal/mapper"
	"github.com/fajarardiyanto/chat-application/internal/model"
	"github.com/fajarardiyanto/chat-application/internal/repo"
	"github.com/fajarardiyanto/chat-application/pkg/auth"
	"github.com/fajarardiyanto/chat-application/pkg/utils"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type AgentProfileController struct {
	accountRepository      repo.AccountRepository
	agentProfileRepository repo.AgentProfileRepository
	agentCredential        repo.AgentCredentialRepository
}

func NewAgentProfileController(
	agentProfileRepository repo.AgentProfileRepository,
	accountRepository repo.AccountRepository,
	agentCredential repo.AgentCredentialRepository,
) *AgentProfileController {
	return &AgentProfileController{
		agentProfileRepository: agentProfileRepository,
		accountRepository:      accountRepository,
		agentCredential:        agentCredential,
	}
}

func (s *AgentProfileController) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	req := request.AgentProfileRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		config.GetLogger().Error(err)
		model.MessageError(w, http.StatusBadRequest, "error decoding request object")
		return
	}

	account, err := s.accountRepository.FindAccountByAccountId(req.AccountId)
	if err != nil {
		config.GetLogger().Error(err)
		model.MessageError(w, http.StatusNotFound, "account not found")
		return
	}

	if _, err = s.agentProfileRepository.FindAgentProfileByEmail(req.Email); err == nil {
		config.GetLogger().Error(err)
		model.MessageError(w, http.StatusFound, "email already exist")
		return
	}

	data := model.AgentProfile{
		Uuid:       uuid.NewString(),
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Email:      req.Email,
		Phone:      req.Phone,
		EmployeeId: req.EmployeeId,
		AccountId:  req.AccountId,
		Role:       req.Role,
		Manager:    req.Manager,
		Audit: model.Audit{
			CreatedAt: time.Now(),
		},
	}

	res, err := s.agentProfileRepository.Register(data)
	if err != nil {
		config.GetLogger().Error(err.Error())
		model.MessageError(w, http.StatusInternalServerError, "something went wrong while registering the user. please try again after sometime.")
		return
	}

	model.MessageSuccess(w, http.StatusOK, mapper.AgentProfileMapper(res, account))
}

func (s *AgentProfileController) SetPassword(w http.ResponseWriter, r *http.Request) {
	req := request.LoginRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		config.GetLogger().Error(err)
		model.MessageError(w, http.StatusBadRequest, "error decoding request object")
		return
	}

	if _, err := s.agentProfileRepository.FindAgentProfileByEmail(req.Email); err != nil {
		config.GetLogger().Error(err)
		model.MessageError(w, http.StatusNotFound, "email not found")
		return
	}

	if _, err := s.agentCredential.FindAgentCredentialByUsername(req.Email); err == nil {
		config.GetLogger().Error(err)
		model.MessageError(w, http.StatusFound, "agent already set password")
		return
	}

	if err := s.agentCredential.SetPassword(req); err != nil {
		config.GetLogger().Error(err)
		model.MessageError(w, http.StatusFound, err.Error())
		return
	}
}

func (s *AgentProfileController) Login(w http.ResponseWriter, r *http.Request) {
	req := request.LoginRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.MessageError(w, http.StatusBadRequest, "error decoding request object")
		return
	}

	agentProfile, err := s.agentProfileRepository.FindAgentProfileByEmail(req.Email)
	if err != nil {
		config.GetLogger().Error(err)
		model.MessageError(w, http.StatusBadRequest, "invalid username/password")
		return
	}

	agentCredential, err := s.agentCredential.FindAgentCredentialByUsername(req.Email)
	if err != nil {
		config.GetLogger().Error(err)
		model.MessageError(w, http.StatusBadRequest, "invalid username/password")
		return
	}

	if err = utils.VerifyPassword(agentCredential.Password, req.Password); err != nil {
		config.GetLogger().Error(err)
		model.MessageError(w, http.StatusBadRequest, "invalid username/password")
		return
	}

	identity := model.TokenModel{
		UserId:    agentProfile.Uuid,
		AccountId: agentProfile.AccountId,
		Role:      constant.AgentRole(agentProfile.Role),
		SessionId: uuid.NewString(),
	}

	token, err := auth.CreateToken(identity)
	if err != nil {
		config.GetLogger().Error(err)
		model.MessageError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	model.MessageSuccess(w, http.StatusOK, mapper.LoginMapper(agentProfile, token))
}
