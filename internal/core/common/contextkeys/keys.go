package contextkeys

type contextKey string

const (
	HTTPRequestKey contextKey = "httpRequest"
	UserIDKey      contextKey = "userID"
)
