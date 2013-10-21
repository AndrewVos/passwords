package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"path"
	"strings"
)

var credentials []Credential
var storedPassword string
var PasswordsFilePath string

func init() {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	PasswordsFilePath = path.Join(usr.HomeDir, ".passwords.passwords")

	fmt.Println(PasswordsFilePath)
	http.HandleFunc("/create_passwords_file", createPasswordsFileHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logged_in", loggedInHandler)
	http.HandleFunc("/passwords_file_exists/", passwordsFileExistsHandler)
	http.HandleFunc("/search/", searchHandler)
	http.HandleFunc("/store", storeHandler)
}

func main() {
	http.ListenAndServe(":8080", nil)
}

func login(password string) bool {
	file, err := ioutil.ReadFile(PasswordsFilePath)
	if err != nil {
		return false
	}
	content, decrypted := Decrypt(file, password)
	if decrypted == false {
		return false
	}
	err = json.Unmarshal(content, &credentials)
	if err != nil {
		return false
	}
	return true
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	response := map[string]interface{}{}
	loggedIn := login(r.FormValue("password"))
	response["logged_in"] = loggedIn

	if loggedIn {
		storedPassword = r.FormValue("password")
	} else {
		storedPassword = ""
	}

	b, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func loggedInHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{}
	response["logged_in"] = storedPassword != ""
	b, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func createPasswordsFileHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	encrypted := Encrypt([]byte("[]"), r.FormValue("password"))
	err := ioutil.WriteFile(PasswordsFilePath, encrypted, 0644)
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
	if _, err := os.Stat(PasswordsFilePath); err == nil {
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
	credentials = append(credentials, credential)
	b, err := json.Marshal(credentials)
	if err != nil {
		panic(err)
	}

	encrypted := Encrypt(b, storedPassword)
	err = ioutil.WriteFile(PasswordsFilePath, encrypted, 0644)
	if err != nil {
		panic(err)
	}
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
