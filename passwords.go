package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var credentials []Credential

func main() {
	http.HandleFunc("/search/", searchHandler)
	http.HandleFunc("/store", storeHandler)
	http.ListenAndServe(":8080", nil)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	for _, credential := range credentials {
		if strings.HasPrefix(credential.Site, query) {
			b, _ := json.Marshal(credential)
			w.Write(b)
			return
		}
	}
}

func storeHandler(w http.ResponseWriter, r *http.Request) {
	credential := Credential{}
	r.ParseForm()

	credential.Site = r.FormValue("site")
	for key, value := range r.Form {
		if fieldNameIsUsername(key) {
			credential.Username = value[0]
		} else if fieldNameIsPassword(key) {
			credential.Password = value[0]
		}
	}
	fmt.Println(credential)
	credentials = append(credentials, credential)
}

func fieldNameIsUsername(field string) bool {
	names := []string{
		"email",
	}
	for _, name := range names {
		if field == name {
			return true
		}
	}
	return false
}

func fieldNameIsPassword(field string) bool {
	names := []string{
		"password",
	}
	for _, name := range names {
		if field == name {
			return true
		}
	}
	return false
}
