package http

import (
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type OAuth interface {
	Login(w http.ResponseWriter, r *http.Request)
	// Callback(handler http.Handler) http.Handler
}

type oauth struct {
	auth oauth2.Config
}

func NewOAuth(
	clientID string,
	clientSecret string,
) OAuth {
	auth := oauth2.Config{
		RedirectURL:  "http://localhost/api/v1/account/google/callback",
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint: google.Endpoint,
	}

	return &oauth{
		auth: auth,
	}
}

func (o *oauth) Login(w http.ResponseWriter, r *http.Request) {
	url := o.auth.AuthCodeURL("state")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
