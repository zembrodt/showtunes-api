package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zembrodt/showtunes-api"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/spotify"
	"log"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"time"

	"github.com/gorilla/mux"
)

type httpResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type pingResponse struct {
	Response httpResponse `json:"response"`
	Name     string       `json:"name"`
	Version  string       `json:"version"`
	ApiRoot  string       `json:"apiRoot"`
}

type httpSuccess httpResponse

type httpError httpResponse

type headerKey int

const (
	contentTypeKey  = "Content-Type"
	contentTypeJSON = "application/json"

	jsonKeyError        = "error"
	jsonKeySuccess      = "success"
	errorInternalServer = "Internal server error"
	errorInvalidRequest = "Invalid request"
	errorInvalidPayload = "Invalid payload"

	pingPath = "/ping"
)

type ShowTunesAPIController struct {
	router        *mux.Router
	routerApi     *mux.Router
	resourcesPath string
	clientId      string
	clientSecret  string
	conf          *oauth2.Config
}

func New(clientId, clientSecret string) *ShowTunesAPIController {
	r := mux.NewRouter()
	rApi := r.PathPrefix(showtunes.APIRoot).Subrouter()
	controller := &ShowTunesAPIController{
		router:       r,
		routerApi:    rApi,
		clientId:     clientId,
		clientSecret: clientSecret,
		conf: &oauth2.Config{
			ClientID: clientId,
			ClientSecret: clientSecret,
			Scopes: scopes,
			Endpoint: spotify.Endpoint,
		},
	}

	// Add handlers
	controller.createGeneralHandlers()
	controller.createAuthHandlers()
	controller.createColorHandlers()

	// Add middlewares for all endpoints
	controller.router.Use(corsMiddleware)
	controller.router.Use(mux.CORSMethodMiddleware(controller.router))
	controller.router.Use(loggerMiddleware)
	controller.router.Use(recoveryMiddleware)

	return controller
}

func (c *ShowTunesAPIController) Start(address string, port int) {
	srv := &http.Server{
		Addr: fmt.Sprintf("%s:%d", address, port),
		WriteTimeout: time.Second * 15,
		ReadTimeout: time.Second * 15,
		IdleTimeout: time.Second * 60,
		Handler: c.router,
	}

	// Run server in a goroutine so it doesn't block
	go func() {
		log.Printf("Starting server on port %d", port)
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	ch := make(chan os.Signal, 1)

	// Graceful shutdown accepted from SIGINT (Ctrl+C)\
	signal.Notify(ch, os.Interrupt)

	// Block until SIGINT
	<-ch

	// Deadline to wait for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 15)
	defer cancel()

	srv.Shutdown(ctx)

	log.Println("Shutting down")
	os.Exit(0)
}

func (c *ShowTunesAPIController) createGeneralHandlers() {
	c.handleGeneralFunc(pingPath, c.ping, http.MethodGet)
}

func (c *ShowTunesAPIController) ping(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, pingResponse{
		Response: httpResponse{
			Code: http.StatusOK,
			Message: jsonKeySuccess,
		},
		Name:    showtunes.Name,
		Version: showtunes.Version,
		ApiRoot: showtunes.APIRoot,
	})
}

// Wrapper to call HandleFunc on the Mux securedRouter and track API endpoints
// Defaults to use all security middleware
// Used for all routes that are prepended with API root
func (c *ShowTunesAPIController) handleFunc(path string, f http.HandlerFunc, methods ...string) {
	c.handleFuncRouter(path, f, c.routerApi, methods...)
}

// Wrapper used for all paths beginning at root
func (c *ShowTunesAPIController) handleGeneralFunc(path string, f http.HandlerFunc, methods ...string) {
	c.handleFuncRouter(path, f, c.router, methods...)
}

// Wrapper for mux.Router.HandleFunc
func (c *ShowTunesAPIController) handleFuncRouter(path string, f http.HandlerFunc, r *mux.Router, methods ...string) {
	methods = append(methods, http.MethodOptions)
	r.HandleFunc(path, f).Methods(methods...)
}

// Write a JSON response with the given HTTP code and payload
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		switch e := err.(type) {
		case *json.UnsupportedTypeError:
			log.Printf("Error: json.Marshal(%v) (err=%v) encountered unsupported type %v, responding with error", payload, e, reflect.TypeOf(payload))
		case *json.UnsupportedValueError:
			log.Printf("Error: json.Marshal(%v) (err=%v) encountered unsupported value of type %v, responding with error", payload, e, reflect.TypeOf(payload))
		default:
			// Shouldn't be able to get here
			log.Printf("Error: json.Marshal(%v) (err=%v) encountered unknown error with payloda of type %v, responding with error", payload, e, reflect.TypeOf(payload))
		}
		respondWithError(w, http.StatusInternalServerError, errorInvalidPayload)
	}

	w.Header().Set(contentTypeKey, contentTypeJSON)
	w.WriteHeader(code)
	w.Write(response)
}

// Wrapper for respondWithJSON to write a success response
func respondWithSuccess(w http.ResponseWriter) {
	respondWithJSON(w, http.StatusOK, httpSuccess{
		Code: http.StatusOK,
		Message: jsonKeySuccess,
	})
}

// Wrapper for respondWithJSON to write an error message
func respondWithError(w http.ResponseWriter, code int, message string, params ...interface{}) {
	respondWithJSON(w, code, httpError{
		Code:    code,
		Message: fmt.Sprintf(message, params...),
	})
}
