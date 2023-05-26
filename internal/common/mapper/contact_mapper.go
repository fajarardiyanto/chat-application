package mapper

import (
	"github.com/fajarardiyanto/chat-application/internal/controller/dto/request"
	"github.com/fajarardiyanto/chat-application/internal/controller/dto/response"
)

func RegisterContactMapper(contact *request.RegisterContactRequest, token string) (res *response.RegisterContactResponse) {
	res = &response.RegisterContactResponse{
		Name:  contact.Name,
		Email: contact.Email,
		Token: token,
	}

	return res
}
