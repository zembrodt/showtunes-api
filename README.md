# ShowTunes API

Go server API as an example of a 3rd party API for the ShowTunes app.

See: https://github.com/zembrodt/showtunes

The ShowTunes app can be configured to use a 3rd party API for its authorization requests. It can also be configured to
use a 3rd party API (the same or a different server) to request analysis information on album covers.

## Endpoints

###`/v1/auth/token` POST
Retrieve/Refresh a Spotify auth token

#### Request
| Parameter       | Value                                                                             |
| --------------- | --------------------------------------------------------------------------------- |
| `grant_type`    | Must be `authorization_token` or `refresh_token`                                  |
| `code`          | The authorization code                                                            |
| `redirect_uri`  | Must match the `redirect_uri` used when requesting the authorization code         |
| `refresh_token` | The token used in place of the authorization code when the auth token has expired |

#### Response
| Parameter       | Type     | Value                                                |
| --------------- | -------- | ---------------------------------------------------- |
| `access_token`  | `string` | The Spotify API access token                         |
| `token_type`    | `string` | How the access token can be used. Always `Bearer`    |
| `refresh_token` | `string` | The token used to request a new token after `expiry` |
| `expiry`        | `string` | Date formatted string for when this token expires    |

---

###`/v1/color` GET
Retrieve the dominant color of a given album's cover art

#### Request
| Parameter | Value                                                                             |
| --------- | --------------------------------------------------------------------------------- |
| `url`     | The url for the image to be used (Must be a domain configured in `VALID_DOMAINS`) |

#### Response
| Parameter | Type     | Value                              |
| --------- | -------- | ---------------------------------- |
| `color`   | `string` | The dominant color as a hex string |

---

###`/ping` GET
Retrieve information on the running API server

#### Response
| Parameter | Type     | Value                             |
| --------- | -------- | ----------------------------------|
| `name`    | `string` | The name of the application       |
| `version`  | `string` | The application's current version |
| `api_root` | `string` | The API endpoint root             |

## Configurations
*Note*: these can be configured as environment variables or in `config/config.yaml`.

Environment variables must be prefixed with `SHOWTUNES_`

| Config           | Default Value | Description                                                        |
| ---------------- | ------------- | ------------------------------------------------------------------ |
| `SERVER_ADDRESS` | `localhost`   | The address of this API server                                     |
| `SERVER_PORT`    | `8000`        | The port of this API server                                        |
| `ORIGIN`         | `*`           | URL for client accessing this API (`Access-Control-Allow-Origin`)  |
| `MAX_AGE`        | `86400`       | Value for `Access-Control-Max-Age` header                          |
| `CLIENT_ID`      | None          | The Client ID to retrieve the authorization token with             |
| `CLIENT_SECRET`  | None          | The Client Secret to retrieve the authorization token with         |
| `VALID_DOMAINS`  | `i.scdn.co`   | Comma-separated list of URLs that host the required Spotify images |

## Building and Running the Server
Scripts have been provided to build an executable for the server that can be deployed.
 * `resources/scripts/build.bat`
 * `resources/scripts/build.sh`
