package handlers

import (
	"Info441FinalProject/servers/gateway/models/users"
	"time"
)

type SessionState struct {
	SessionStartTime  time.Time
	AuthenticatedUser *users.User
}
