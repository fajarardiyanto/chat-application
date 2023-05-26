package repository

import "github.com/fajarardiyanto/chat-application/internal/model"

type ContactRepository interface {
	RegisterContact(contact model.Contact) error
	FindByEmailAndAccountUuid(email string, accountId string) (*model.Contact, error)
}
