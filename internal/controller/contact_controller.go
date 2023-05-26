package controller

import (
	"encoding/json"
	"github.com/fajarardiyanto/chat-application/config"
	"github.com/fajarardiyanto/chat-application/internal/common/exception"
	"github.com/fajarardiyanto/chat-application/internal/common/mapper"
	"github.com/fajarardiyanto/chat-application/internal/controller/dto/request"
	"github.com/fajarardiyanto/chat-application/internal/model"
	"github.com/fajarardiyanto/chat-application/internal/repository"
	"github.com/fajarardiyanto/chat-application/pkg/auth"
	"github.com/fajarardiyanto/chat-application/pkg/utils"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type ContactController struct {
	contactRepository          repository.ContactRepository
	contactInboxRepository     repository.ContactInboxRepository
	channelWebWidgetRepository repository.ChannelWebWidgetRepository
	inboxRepository            repository.InboxRepository
}

func NewContactController(
	contactRepository repository.ContactRepository,
	contactInboxRepository repository.ContactInboxRepository,
	channelWebWidgetRepository repository.ChannelWebWidgetRepository,
	inboxRepository repository.InboxRepository,
) *ContactController {
	return &ContactController{
		contactRepository:          contactRepository,
		contactInboxRepository:     contactInboxRepository,
		channelWebWidgetRepository: channelWebWidgetRepository,
		inboxRepository:            inboxRepository,
	}
}

func (s *ContactController) RegisterContact(w http.ResponseWriter, r *http.Request) {
	websiteToken := utils.QueryParam(r, "websiteToken")

	req := &request.RegisterContactRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		config.GetLogger().Error(err)
		model.MessageError(w, http.StatusBadRequest, exception.ErrorDecodeRequest)
		return
	}

	channel, err := s.channelWebWidgetRepository.FindByWebsiteToken(websiteToken)
	if err != nil {
		config.GetLogger().Error(err.Error())
		model.MessageError(w, http.StatusNotFound, exception.WebsiteTokenNotFound)
		return
	}

	inbox, err := s.inboxRepository.FindByChannelId(channel.Uuid)
	if err != nil {
		config.GetLogger().Error(err.Error())
		model.MessageError(w, http.StatusNotFound, exception.ChannelWithInboxNotFound)
		return
	}

	contactId := uuid.NewString()
	contact, err := s.contactRepository.FindByEmailAndAccountUuid(req.Email, inbox.AccountId)
	if err != nil {
		contactData := model.Contact{
			Name:      req.Name,
			Email:     req.Email,
			Phone:     req.Phone,
			Uuid:      contactId,
			AccountId: inbox.AccountId,
			CreatedAt: time.Now(),
		}

		if err = s.contactRepository.RegisterContact(contactData); err != nil {
			config.GetLogger().Error(err.Error())
			model.MessageError(w, http.StatusInternalServerError, exception.SomethingWentWrong)
			return
		}
	} else {
		contactId = contact.Uuid
	}

	sourceId := uuid.NewString()
	contactInbox := model.ContactInbox{
		ContactId:   contactId,
		InboxId:     inbox.Uuid,
		Uuid:        uuid.NewString(),
		PubSubToken: uuid.NewString(),
		SourceId:    sourceId,
		CreatedAt:   time.Now(),
	}

	if err = s.contactInboxRepository.CreateContactInbox(contactInbox); err != nil {
		config.GetLogger().Error(err.Error())
		model.MessageError(w, http.StatusInternalServerError, exception.SomethingWentWrong)
		return
	}

	identity := model.ContactTokenModel{
		SourceId: sourceId,
		InboxId:  inbox.Uuid,
	}

	token, err := auth.CreateTokenContact(identity)
	if err != nil {
		config.GetLogger().Error(err)
		model.MessageError(w, http.StatusInternalServerError, err.Error())
		return
	}

	model.MessageSuccess(w, http.StatusOK, mapper.RegisterContactMapper(req, token))
}
