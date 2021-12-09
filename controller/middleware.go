package controller

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"github.com/zembrodt/showtunes-api/util/global"
	"log"
	"net/http"
	"strings"
)

func allowedHeaders() string {
	return strings.Join([]string{
		"Content-Type", "accept", "origin",
	}, ", ")
}

func allowedMethods() string {
	return strings.Join([]string{
		http.MethodOptions, http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete,
	}, ", ")
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", viper.GetString(global.OriginKey))
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Allow-Headers", allowedHeaders())
		w.Header().Add("Access-Control-Allow-Methods", allowedMethods())
		w.Header().Add("Access-Control-Max-Age", viper.GetString(global.MaxAgeKey))

		if r.Method == http.MethodOptions {
			http.Error(w, "No Content", http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		log.Printf("%s [%s]", r.URL, r.Method)
		next.ServeHTTP(w, r)
	})
}

func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			var err error
			rec := recover()
			if rec != nil {
				switch t := rec.(type) {
				case string:
					err = errors.New(t)
				case error:
					err = t
				default:
					err = fmt.Errorf("unexpected error type %v: %v", t, rec)
				}
				log.Printf("Recovered from panic: %v", err)
				http.Error(w,"Unexpected error occurred", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
