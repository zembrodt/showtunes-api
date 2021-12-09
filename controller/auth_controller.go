package controller

import (
	"context"
	"golang.org/x/oauth2"
	"net/http"
)

const (
	pathAuth = "/auth"
	pathToken = pathAuth + "/tokens"

	codeKey          = "code"
	redirectUriKey   = "redirect_uri"
)

var scopes = []string{
	"user-read-playback-state",
	"user-modify-playback-state",
}

func (c *ShowTunesAPIController) createAuthHandlers() {
	c.handleFunc(pathToken, c.getAuthTokens, http.MethodPost)
	c.handleFunc(pathToken, c.updateAuthToken, http.MethodPut)
}

func (c *ShowTunesAPIController) getAuthTokens(w http.ResponseWriter, r *http.Request) {
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

func (c *ShowTunesAPIController) updateAuthToken(w http.ResponseWriter, r *http.Request) {
	// Get request parameters
	code := r.FormValue(codeKey)

	if len(code) == 0 {
		respondWithError(w, http.StatusBadRequest, "Invalid %s", codeKey)
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
