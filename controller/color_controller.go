package controller

import (
	"bytes"
	"github.com/cenkalti/dominantcolor"
	"github.com/zembrodt/music-display-api/model"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	pathColor = "/color"

	urlKey = "url"

	barcodeColorThreshold = 186
)

var validDomains = map[string]bool {
	"i.scdn.co": true,
}

func (c *MusicAPIController) createColorHandlers() {
	c.handleFunc(pathColor, c.getDominateColor, http.MethodGet)
}

func (c *MusicAPIController) getDominateColor(w http.ResponseWriter, r *http.Request) {
	urls, ok := r.URL.Query()[urlKey]
	if !ok || len(urls[0]) < 1 {
		respondWithError(w, http.StatusBadRequest, "Request does not contain id")
		return
	}

	encodedUrl := urls[0]
	// TODO: get image from url

	coverArtUrlRaw, err := url.QueryUnescape(encodedUrl)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to unescape the provided cover art url: " + err.Error())
		return
	}

	coverArtUrl, err := url.Parse(coverArtUrlRaw)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to parse valid URL: " + err.Error())
		return
	}
	// Check if the domain for this image is an accepted one
	if !validDomains[strings.ToLower(coverArtUrl.Hostname())] {
		respondWithError(w, http.StatusBadRequest, "The provided domain in the URL is invalid")
		return
	}
	// Check if this url uses https
	if strings.ToLower(coverArtUrl.Scheme) != "https" {
		respondWithError(w, http.StatusBadRequest, "Provided URL must use HTTPS")
		return
	}

	res, err := http.Get(coverArtUrlRaw)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Provided image URL was not successfully retrieved: " + err.Error())
		return
	}
	defer res.Body.Close()

	imageData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to retrieve data from provided URL: " + err.Error())
		return
	}

	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Provided URL is not an image or of an unsupported type: " + err.Error())
		return
	}

	// Get dominant color of the album art
	var colorResponse model.ColorResponse
	rgb := dominantcolor.Find(img)
	colorResponse.Color = dominantcolor.Hex(rgb)
	// See https://stackoverflow.com/questions/3942878/how-to-decide-font-color-in-white-or-black-depending-on-background-color
	colorResponse.UseBlackBarcodeColor = float32(rgb.R) * 0.299 + float32(rgb.G) * 0.587 + float32(rgb.B) * 0.114 > barcodeColorThreshold

	respondWithJSON(w, http.StatusOK, colorResponse)
}
