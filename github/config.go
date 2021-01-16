package github

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
)

var oauthConfig = &oauth2.Config{
	ClientID:     "bb4d363634444833a6a0",
	ClientSecret: "419adb405f31547fe13436e5daa0c6c836b654bc",
	Endpoint:     endpoints.GitHub,
	Scopes: []string{
		"user:email",
	},
}
