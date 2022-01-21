package pa

/*
	You could argue that the session types should be under the http package but I consider them important
	enough to get their own package + you could implement a session cacher which would also need access
	to the session type
*/

// SessionCookieName represents the name of the session cookie.
const SessionCookieName = "session"

// Session represents data stored per session under a secure cookie.
type Session struct {
	UserID  int    `json:"userID"`
	IsAdmin bool   `json:"isAdmin"`
	State   string `json:"state"`
}
