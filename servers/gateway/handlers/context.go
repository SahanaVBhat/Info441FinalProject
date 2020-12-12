package handlers

import (
	"github.com/info441-au20/assignments-melodc/servers/gateway/models/users"
	"github.com/info441-au20/assignments-melodc/servers/gateway/sessions"
)

type HandlerContext struct {
	SigningKey   string
	SessionStore sessions.Store
	UserStore    users.Store
}
