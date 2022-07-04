package utils

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// HandleErrorCode Generates error page
func HandleErrorCode(errorCode int, description string, w http.ResponseWriter) {
	w.WriteHeader(errorCode)                    // set HTTP status code (example 404, 500)
	w.Header().Set("Content-Type", "text/html") // clarify return type (MIME)
	w.Write([]byte(fmt.Sprintf(
		"<html><body><h1>Error %d</h1><p>%s</p></body></html>",
		errorCode,
		description)))
}

func HandleErrorCodeLog(errorCode int, err error, w http.ResponseWriter) {
	log.Error(err)
	w.WriteHeader(errorCode)                    // set HTTP status code (example 404, 500)
	w.Header().Set("Content-Type", "text/html") // clarify return type (MIME)
	w.Write([]byte(fmt.Sprintf(
		"<html><body><h1>Error %d</h1><p>%s</p></body></html>",
		errorCode,
		err)))
}
