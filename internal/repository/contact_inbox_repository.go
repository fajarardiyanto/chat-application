package repository

import "github.com/fajarardiyanto/chat-application/internal/model"

type ContactInboxRepository interface {
	CreateContactInbox(contactInbox model.ContactInbox) error
}
