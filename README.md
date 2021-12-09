# ShowTunes API

Go server used by the ShowTunes app to make API requests where a client secret is needed.

## Requests
* */v1/auth/tokens* (**POST**): Retrieve a Spotify auth token
  * Expects *code* and *redirect_uri*
  * Returns a new auth token
* */v1/auth/tokens* (**PUT**): Refresh a Spotify auth token
  * Expects refresh token as *code*
  * Returns an updated auth token
* */v1/color* (**GET**): Retrieve the dominant color of a given album's cover art
  * Expects *url* - must be to an image hosted by Spotify (*i.scdn.co*)
  * Returns the dominant color in hex
