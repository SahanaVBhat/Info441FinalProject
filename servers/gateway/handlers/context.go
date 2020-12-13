package handlers

import (
	"github.com/SahanaVBhat/Info441FinalProject/servers/gateway/models/users"
	"github.com/SahanaVBhat/Info441FinalProject/servers/gateway/sessions"
)

type HandlerContext struct {
	SigningKey   string
	SessionStore sessions.Store
	UserStore    users.Store
}
