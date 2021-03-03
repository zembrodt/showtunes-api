package controller

import (
	"context"
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

func (c *MusicAPIController) createAuthHandlers() {
	c.handleFunc(pathToken, c.getAuthTokens, http.MethodPost)
}

func (c *MusicAPIController) getAuthTokens(w http.ResponseWriter, r *http.Request) {
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
