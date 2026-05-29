package session

import "time"

type Session struct {
	Id                 string
	IPAdress           string
	UserAgent          string
	Username           string
	Email              string
	LanguagePreference string
	CurrentState       string
	LastLogin          time.Time
	TTL                time.Time
}

// evicted by TTL
func NewPersistentSession(Id string,
	IPAdress string,
	UserAgent string,
	Username string,
	Email string,
	LanguagePreference string,
	CurrentState string,
	LastLogin time.Time,
	TTL time.Time) *Session {

	return &Session{
		Id:                 Id,
		IPAdress:           IPAdress,
		UserAgent:          UserAgent,
		Username:           Username,
		Email:              Email,
		LanguagePreference: LanguagePreference,
		CurrentState:       CurrentState,
		LastLogin:          LastLogin,
		TTL:                TTL,
	}

}

// evicted manually
func NewNonPersistentSession(Id string,
	IPAdress string,
	UserAgent string,
	Username string,
	Email string,
	LanguagePreference string,
	CurrentState string,
	LastLogin time.Time,
) *Session {

	return &Session{
		Id:                 Id,
		IPAdress:           IPAdress,
		UserAgent:          UserAgent,
		Username:           Username,
		Email:              Email,
		LanguagePreference: LanguagePreference,
		CurrentState:       CurrentState,
		LastLogin:          LastLogin,
	}

}
