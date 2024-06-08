// This package includes web api
package web

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type EndpointHandler func(w http.ResponseWriter, req *http.Request)

func HealthCheckEndpoint(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Keep calm I'm absolutely alive 🐛")
}

func NoRouteEndpoint(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "Oops! 👀")
}

func HandleErrors(w http.ResponseWriter, err error, httpStatus int) {
	w.WriteHeader(httpStatus)
	log.Errorf("⛔ %v", err.Error())
	fmt.Fprintf(w, "%v", err.Error())
}
