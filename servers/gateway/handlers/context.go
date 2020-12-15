package handlers

import (
	"Info441FinalProject/servers/gateway/models/users"
	"Info441FinalProject/servers/gateway/sessions"
)

type HandlerContext struct {
	SigningKey   string
	SessionStore sessions.Store
	UserStore    users.Store
}
