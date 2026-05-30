package session

import "time"

type Session struct {
	Id         string
	UserID     string
	CreatedAt  time.Time
	LastAccess time.Time
	TTL        time.Time
}

// evicted by TTL
func NewSession(
	Id string,
	UserID string,
	CreatedAt time.Time,
	LastAccess time.Time,
	TTL time.Time) *Session {

	return &Session{
		Id:         Id,
		UserID:     UserID,
		CreatedAt:  CreatedAt,
		LastAccess: LastAccess,
		TTL:        TTL,
	}

}
