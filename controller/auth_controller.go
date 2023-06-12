package controller

import (
	"context"
	"golang.org/x/oauth2"
	"net/http"
)

const (
	pathAuth = "/auth"
	pathToken = pathAuth + "/token"

	codeKey         = "code"
	redirectUriKey  = "redirect_uri"
	grantTypeKey    = "grant_type"
	refreshTokenKey = "refresh_token"
)

func (c *ShowTunesAPIController) createAuthHandlers() {
	c.handleFunc(pathToken, c.handleToken, http.MethodPost)
}

func (c *ShowTunesAPIController) handleToken(w http.ResponseWriter, r *http.Request) {
	if r.FormValue(grantTypeKey) == refreshTokenKey {
		c.refreshToken(w, r)
	} else {
		c.fetchToken(w, r)
	}
}

func (c *ShowTunesAPIController) fetchToken(w http.ResponseWriter, r *http.Request) {
	// Get request parameters
	code := r.FormValue(codeKey)
	redirectUri := r.FormValue(redirectUriKey)

	if len(code) == 0 {
		respondWithError(w, http.StatusBadRequest, "Invalid %s", codeKey)
		return
	}
	if len(redirectUri) == 0 {
		respondWithError(w, http.StatusBadRequest, "Invalid %s", redirectUriKey)
		return
	}

	// Send token request to Spotify
	c.conf.RedirectURL = redirectUri
	token, err := c.conf.Exchange(context.Background(), code)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Invalid request to Spotify token endpoint: %s", err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, token)
}

func (c *ShowTunesAPIController) refreshToken(w http.ResponseWriter, r *http.Request) {
	// Get request parameters
	code := r.FormValue(refreshTokenKey)

	if len(code) == 0 {
		respondWithError(w, http.StatusBadRequest, "Invalid %s", refreshTokenKey)
		return
	}

	refreshToken := &oauth2.Token{
		RefreshToken: code,
	}

	tokenSource := c.conf.TokenSource(context.Background(), refreshToken)
	newToken, err := tokenSource.Token()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to refresh auth token: %s", err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, newToken)
}
