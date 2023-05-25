package exception

const (
	AgentNotFound           = "agent not found"
	ErrorDecodeRequest      = "error decoding request object"
	AccountNotFound         = "account not found"
	EmailAlreadyExist       = "email already exist"
	FailedRegister          = "something went wrong while registering the user. please try again after sometime"
	NotAllowedToSetPassword = "you are not allowed to set password"
	EmailNotFound           = "email not found"
	AgentAlreadySetPassword = "agent already set password"
	InvalidUsernamePassword = "invalid username/password"
	InvalidToken            = "invalid token"
	NotAllowed              = "you are not allowed!"
	NotAllowedToChat        = "this role is not in allowed list to chat with customer"
	ConversationNotFound    = "conversation not found"
)