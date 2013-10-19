package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var loadedCredentials = false
var credentials []Credential

func main() {
	http.HandleFunc("/create_passwords_file", createPasswordsFileHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/passwords_file_exists/", passwordsFileExistsHandler)
	http.HandleFunc("/search/", searchHandler)
	http.HandleFunc("/store", storeHandler)
	http.ListenAndServe(":8080", nil)
}

func login(password string) ([]Credential, bool) {
	file, err := ioutil.ReadFile("passwords_file")
	if err != nil {
		return nil, false
	}
	content, decrypted := Decrypt(file, password)
	if decrypted == false {
		return nil, false
	}
	var credentials []Credential
	json.Unmarshal(content, &credentials)
	return credentials, true
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	response := map[string]interface{}{}
	c, loggedIn := login(r.FormValue("password"))
	credentials = c
	response["logged_in"] = loggedIn
	b, _ := json.Marshal(response)
	w.Write(b)
}

func createPasswordsFileHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	encrypted := Encrypt([]byte("[]"), r.FormValue("password"))
	err := ioutil.WriteFile("passwords_file", encrypted, 0777)
	if err != nil {
		panic(err)
	}
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

func passwordsFileExistsHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{}
	if _, err := os.Stat("passwords_file"); err == nil {
		response["passwords_file_exists"] = true
	} else {
		response["passwords_file_exists"] = false
	}
	b, _ := json.Marshal(response)
	w.Write(b)
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
