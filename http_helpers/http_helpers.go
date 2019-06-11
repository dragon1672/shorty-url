package http_helpers

import (
	"encoding/json"
	"github.com/golang/glog"
	"net/http"
	"net/url"
)

func Redirect(url string, w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func PrintError(code int, message string, w http.ResponseWriter) {
	err := writeResponse(code, map[string]string{"error": message}, w)
	if err != nil {
		glog.Infof("error when writing response: %v", err)
	}
}

func PrintText(message string, w http.ResponseWriter) {
	_, err := w.Write([]byte(message))
	if err != nil {
		glog.Infof("error when writing response: %v", err)
	}
}

func writeResponse(code int, payload interface{}, w http.ResponseWriter) error {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err := w.Write(response)
	return err
}

func IsValidURL(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	return err == nil
}
