package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

func setup() {
	removePasswordsFile()
}

func teardown() {
	removePasswordsFile()
}

func removePasswordsFile() {
	os.Remove("passwords_file")
}

func getJSON(server *httptest.Server, path string) map[string]interface{} {
	var responseData map[string]interface{}
	response, _ := http.Get(server.URL + "/passwords_file_exists")
	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(data, &responseData)
	response.Body.Close()
	return responseData
}

func postFormJSON(server *httptest.Server, path string, form url.Values) {
	http.PostForm(server.URL+path, form)
}

func TestPasswordsFileExists(t *testing.T) {
	setup()
	ts := httptest.NewServer(http.HandlerFunc(passwordsFileExistsHandler))

	responseData := getJSON(ts, "/passwords_file_exists")
	if responseData["passwords_file_exists"] != false {
		t.Errorf("Expected passwords_file_exists = false")
	}

	os.Create("passwords_file")
	responseData = getJSON(ts, "/passwords_file_exists")
	if responseData["passwords_file_exists"] != true {
		t.Errorf("Expected passwords_file_exists = true")
	}
	teardown()
}

func TestCreatePasswordsFile(t *testing.T) {
	setup()
	ts := httptest.NewServer(http.HandlerFunc(createPasswordsFileHandler))

	form := url.Values{"password": {"some-password"}}
	postFormJSON(ts, "/create_passwords_file", form)
	if _, err := os.Stat("passwords_file"); err != nil {
		t.Errorf("Expected passwords_file to be created")
	}
	teardown()
}

func TestLogin(t *testing.T) {
	// autofill
	// check passwordsfile exists
	// if it doesn't create account
	// if it does, checked logged in state
}

func TestStoreAndSearchPassword(t *testing.T) {
	// ts := httptest.NewServer(http.HandlerFunc(createPasswordsFileHandler))
	// body := bytes.NewReader([]byte(`{"password": "master password"}`))
	// http.Post(ts.URL+"/create_passwords_file", "application/json", body)

	// ts = httptest.NewServer(http.HandlerFunc(storeHandler))
	// body = bytes.NewReader([]byte(`{"site": "http://example.com", "username": "andrew", "password": "a-password"}`))

	// form := url.Values{"site": {"http://example.com"}, "email": {"andrew@me.com"}, "password": {"a-passwod"}}
	// http.PostForm(ts.URL+"/store/", form)

	// ts = httptest.NewServer(http.HandlerFunc(searchHandler))
	// response, _ := http.Get(ts.URL + "/search/?q=" + url.QueryEscape("http://example.com"))
	// defer response.Body.Close()
	// data, _ := ioutil.ReadAll(response.Body)
	// var v []map[string]interface{}
	// json.Unmarshal(data, &v)
	// fmt.Println(v)
}

// login
// search passwords
// edit password
// add passwords
