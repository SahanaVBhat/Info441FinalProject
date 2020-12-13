package handlers

import (
	"time"

	"github.com/SahanaVBhat/Info441FinalProject/servers/gateway/models/users"
)

type SessionState struct {
	SessionStartTime  time.Time
	AuthenticatedUser *users.User
}
