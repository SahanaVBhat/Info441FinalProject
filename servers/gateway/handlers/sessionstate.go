package handlers

import (
	"time"

	"github.com/info441-au20/assignments-melodc/servers/gateway/models/users"
)

type SessionState struct {
	SessionStartTime  time.Time
	AuthenticatedUser *users.User
}
