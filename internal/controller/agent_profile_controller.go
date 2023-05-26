package controller

import (
	"encoding/json"
	"github.com/fajarardiyanto/chat-application/config"
	"github.com/fajarardiyanto/chat-application/internal/common/constant"
	"github.com/fajarardiyanto/chat-application/internal/common/exception"
	"github.com/fajarardiyanto/chat-application/internal/common/mapper"
	"github.com/fajarardiyanto/chat-application/internal/common/validation"
	"github.com/fajarardiyanto/chat-application/internal/controller/dto/request"
	"github.com/fajarardiyanto/chat-application/internal/model"
	"github.com/fajarardiyanto/chat-application/internal/repository"
	"github.com/fajarardiyanto/chat-application/pkg/auth"
	"github.com/fajarardiyanto/chat-application/pkg/utils"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type AgentProfileController struct {
	accountRepository         repository.AccountRepository
	agentProfileRepository    repository.AgentProfileRepository
	agentCredentialRepository repository.AgentCredentialRepository
}

func NewAgentProfileController(
	agentProfileRepository repository.AgentProfileRepository,
	accountRepository repository.AccountRepository,
	agentCredentialRepository repository.AgentCredentialRepository,
) *AgentProfileController {
	return &AgentProfileController{
		agentProfileRepository:    agentProfileRepository,
		accountRepository:         accountRepository,
		agentCredentialRepository: agentCredentialRepository,
	}
}

func (s *AgentProfileController) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	config.GetLogger().Info("starting register agent")

	req := request.AgentProfileRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		config.GetLogger().Error(err)
		model.MessageError(w, http.StatusBadRequest, exception.ErrorDecodeRequest)
		return
	}

	account, err := s.accountRepository.FindAccountByAccountId(req.AccountId)
	if err != nil {
		config.GetLogger().Error(err)
		model.MessageError(w, http.StatusNotFound, exception.AccountNotFound)
		return
	}

	if _, err = s.agentProfileRepository.FindAgentProfileByEmail(req.Email); err == nil {
		config.GetLogger().Error(err)
		model.MessageError(w, http.StatusFound, exception.EmailAlreadyExist)
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
		model.MessageError(w, http.StatusInternalServerError, exception.FailedRegister)
		return
	}

	model.MessageSuccess(w, http.StatusOK, mapper.AgentProfileMapper(res, account))
}

func (s *AgentProfileController) SetPassword(w http.ResponseWriter, r *http.Request) {
	config.GetLogger().Info("starting set password agent")

	if !validation.IsAllowedToSetPassword(r) {
		config.GetLogger().Error(exception.NotAllowedToSetPassword)
		model.MessageError(w, http.StatusUnauthorized, exception.NotAllowedToSetPassword)
		return
	}

	req := request.LoginRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		config.GetLogger().Error(err)
		model.MessageError(w, http.StatusBadRequest, exception.ErrorDecodeRequest)
		return
	}

	if _, err := s.agentProfileRepository.FindAgentProfileByEmail(req.Email); err != nil {
		config.GetLogger().Error(err)
		model.MessageError(w, http.StatusNotFound, exception.EmailNotFound)
		return
	}

	if _, err := s.agentCredentialRepository.FindAgentCredentialByUsername(req.Email); err == nil {
		config.GetLogger().Error(err)
		model.MessageError(w, http.StatusFound, exception.AgentAlreadySetPassword)
		return
	}

	if err := s.agentCredentialRepository.SetPassword(req); err != nil {
		config.GetLogger().Error(err)
		model.MessageError(w, http.StatusFound, err.Error())
		return
	}
}

func (s *AgentProfileController) Login(w http.ResponseWriter, r *http.Request) {
	req := request.LoginRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.MessageError(w, http.StatusBadRequest, exception.ErrorDecodeRequest)
		return
	}

	agentProfile, err := s.agentProfileRepository.FindAgentProfileByEmail(req.Email)
	if err != nil {
		config.GetLogger().Error(err)
		model.MessageError(w, http.StatusBadRequest, exception.AgentNotFound)
		return
	}

	agentCredential, err := s.agentCredentialRepository.FindAgentCredentialByUsername(req.Email)
	if err != nil {
		config.GetLogger().Error(err)
		model.MessageError(w, http.StatusBadRequest, exception.InvalidUsernamePassword)
		return
	}

	if err = utils.VerifyPassword(agentCredential.Password, req.Password); err != nil {
		config.GetLogger().Error(err)
		model.MessageError(w, http.StatusBadRequest, exception.InvalidUsernamePassword)
		return
	}

	identity := model.AgentTokenModel{
		UserId:    agentProfile.Uuid,
		AccountId: agentProfile.AccountId,
		Role:      constant.AgentRole(agentProfile.Role),
		SessionId: uuid.NewString(),
	}

	token, err := auth.CreateTokenAgent(identity)
	if err != nil {
		config.GetLogger().Error(err)
		model.MessageError(w, http.StatusInternalServerError, err.Error())
		return
	}

	model.MessageSuccess(w, http.StatusOK, mapper.LoginMapper(agentProfile, token))
}

func (s *AgentProfileController) Logout(w http.ResponseWriter, r *http.Request) {
	token, err := auth.ExtractTokenAgent(r)
	if err != nil {
		config.GetLogger().Error(err)
		model.MessageError(w, http.StatusBadRequest, exception.InvalidToken)
		return
	}

	_, err = s.agentProfileRepository.FindAgentProfileById(token.UserId)
	if err != nil {
		config.GetLogger().Error(err)
		model.MessageError(w, http.StatusNotFound, exception.AgentNotFound)
		return
	}

}

func (s *AgentProfileController) GetAgent(w http.ResponseWriter, r *http.Request) {
	agentId := utils.QueryString(r, "agentId")

	if !validation.IsAllowedToGetAgent(r) {
		config.GetLogger().Error(exception.NotAllowedToSetPassword)
		model.MessageError(w, http.StatusUnauthorized, exception.NotAllowedToSetPassword)
		return
	}

	token, err := auth.ExtractTokenAgent(r)
	if err != nil {
		config.GetLogger().Error(err)
		model.MessageError(w, http.StatusBadRequest, exception.InvalidToken)
		return
	}

	account, err := s.accountRepository.FindAccountByAccountId(token.AccountId)
	if err != nil {
		config.GetLogger().Error(err)
		model.MessageError(w, http.StatusNotFound, exception.AccountNotFound)
		return
	}

	agentProfile, err := s.agentProfileRepository.FindAgentProfileById(agentId)
	if err != nil {
		config.GetLogger().Error(err)
		model.MessageError(w, http.StatusNotFound, exception.AgentNotFound)
		return
	}

	if agentProfile.AccountId != token.AccountId {
		config.GetLogger().Error(exception.NotAllowed)
		model.MessageError(w, http.StatusUnauthorized, exception.NotAllowed)
		return
	}

	model.MessageSuccess(w, http.StatusOK, mapper.AgentProfileMapper(agentProfile, account))
}

func (s *AgentProfileController) GetAllAgentByAccountId(w http.ResponseWriter, r *http.Request) {
	accountId := utils.QueryParam(r, "accountId")

	if !validation.IsAllowedToGetAllAgent(r) {
		config.GetLogger().Error(exception.NotAllowedToSetPassword)
		model.MessageError(w, http.StatusUnauthorized, exception.NotAllowedToSetPassword)
		return
	}

	token, err := auth.ExtractTokenAgent(r)
	if err != nil {
		config.GetLogger().Error(err)
		model.MessageError(w, http.StatusBadRequest, exception.InvalidToken)
		return
	}

	account, err := s.accountRepository.FindAccountByAccountId(token.AccountId)
	if err != nil {
		config.GetLogger().Error(err)
		model.MessageError(w, http.StatusNotFound, exception.AccountNotFound)
		return
	}

	agentProfile, err := s.agentProfileRepository.FindAgentProfileByAccountId(accountId)
	if err != nil {
		config.GetLogger().Error(err)
		model.MessageError(w, http.StatusNotFound, exception.AgentNotFound)
		return
	}

	model.MessageSuccess(w, http.StatusOK, mapper.AllAgentProfileMapper(agentProfile, account))
}

func (s *AgentProfileController) UpdateAgentProfileById(w http.ResponseWriter, r *http.Request) {
	agentId := utils.QueryParam(r, "agentId")

	req := request.AgentProfileRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.MessageError(w, http.StatusBadRequest, exception.ErrorDecodeRequest)
		return
	}

	token, err := auth.ExtractTokenAgent(r)
	if err != nil {
		config.GetLogger().Error(err)
		model.MessageError(w, http.StatusBadRequest, exception.InvalidToken)
		return
	}

	if token.UserId != agentId {
		config.GetLogger().Error(exception.NotAllowed)
		model.MessageError(w, http.StatusUnauthorized, exception.NotAllowed)
		return
	}

	data := model.AgentProfile{
		Uuid:      agentId,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	if err = s.agentProfileRepository.UpdateAgentProfileById(data); err != nil {
		config.GetLogger().Error(err)
		model.MessageError(w, http.StatusUnauthorized, err.Error())
		return
	}
}

func (s *AgentProfileController) DeleteAgentProfileById(w http.ResponseWriter, r *http.Request) {
	agentId := utils.QueryParam(r, "agentId")

	if !validation.IsAllowedToDeleteAgent(r) {
		config.GetLogger().Error(exception.NotAllowedToSetPassword)
		model.MessageError(w, http.StatusUnauthorized, exception.NotAllowedToSetPassword)
		return
	}

	if err := s.agentProfileRepository.DeleteAgentProfileById(agentId); err != nil {
		config.GetLogger().Error(err)
		model.MessageError(w, http.StatusUnauthorized, err.Error())
		return
	}
}
