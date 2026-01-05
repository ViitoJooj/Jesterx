package config

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

var (
	GoogleOAuthConfig  *oauth2.Config
	GithubOAuthConfig  *oauth2.Config
	TwitterOAuthConfig *oauth2.Config
)

func InitOAuth() {
	if GoogleClientID != "" && GoogleClientSecret != "" {
		GoogleOAuthConfig = &oauth2.Config{
			ClientID:     GoogleClientID,
			ClientSecret: GoogleClientSecret,
			RedirectURL:  GoogleRedirectURL,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		}
	}

	if GithubClientID != "" && GithubClientSecret != "" {
		GithubOAuthConfig = &oauth2.Config{
			ClientID:     GithubClientID,
			ClientSecret: GithubClientSecret,
			RedirectURL:  GithubRedirectURL,
			Scopes:       []string{"user:email"},
			Endpoint:     github.Endpoint,
		}
	}

	if TwitterClientID != "" && TwitterClientSecret != "" {
		TwitterOAuthConfig = &oauth2.Config{
			ClientID:     TwitterClientID,
			ClientSecret: TwitterClientSecret,
			RedirectURL:  TwitterRedirectURL,
			Scopes:       []string{"tweet.read", "users.read"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://twitter.com/i/oauth2/authorize",
				TokenURL: "https://api.twitter.com/2/oauth2/token",
			},
		}
	}
}
