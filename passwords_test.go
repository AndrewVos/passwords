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

func postFormJSON(server *httptest.Server, path string, form url.Values) map[string]interface{} {
	var responseData map[string]interface{}
	response, _ := http.PostForm(server.URL+path, form)
	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(data, &responseData)
	response.Body.Close()
	return responseData
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
	setup()

	ts := httptest.NewServer(http.HandlerFunc(loginHandler))
	form := url.Values{"password": {"some-password"}}
	responseData := postFormJSON(ts, "/login", form)
	if responseData["logged_in"] == true {
		t.Errorf("Expected to not be logged in yet")
	}

	ts = httptest.NewServer(http.HandlerFunc(createPasswordsFileHandler))
	form = url.Values{"password": {"some-password"}}
	postFormJSON(ts, "/create_passwords_file", form)

	ts = httptest.NewServer(http.HandlerFunc(loginHandler))
	form = url.Values{"password": {"some-password"}}
	responseData = postFormJSON(ts, "/login", form)
	if responseData["logged_in"] == false {
		t.Errorf("Expected to be logged in")
	}

	teardown()
}
