package github

import "time"

// random-uuid -> expiration time
var loginAttempts = map[string]time.Time{}

// provider user id -> internal user id
var oauthConnections = map[string]string{}

// sessionID -> internal user id
var sessions = map[string]string{}
