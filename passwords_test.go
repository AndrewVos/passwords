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

var server *httptest.Server

func setup() {
	removePasswordsFile()
	storedPassword = ""
	credentials = nil
	server = httptest.NewServer(nil)
}

func teardown() {
	removePasswordsFile()
}

func removePasswordsFile() {
	os.Remove("passwords_file")
}

func getJSON(path string) map[string]interface{} {
	var responseData map[string]interface{}
	response, _ := http.Get(server.URL + path)
	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(data, &responseData)
	response.Body.Close()
	return responseData
}

func postFormJSON(path string, form url.Values) map[string]interface{} {
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

	responseData := getJSON("/passwords_file_exists")
	if responseData["passwords_file_exists"] != false {
		t.Errorf("Expected passwords_file_exists = false")
	}

	os.Create("passwords_file")
	responseData = getJSON("/passwords_file_exists")
	if responseData["passwords_file_exists"] != true {
		t.Errorf("Expected passwords_file_exists = true")
	}
	teardown()
}

func TestCreatePasswordsFile(t *testing.T) {
	setup()

	form := url.Values{"password": {"some-password"}}
	postFormJSON("/create_passwords_file", form)
	if _, err := os.Stat("passwords_file"); err != nil {
		t.Errorf("Expected passwords_file to be created")
	}
	teardown()
}

func TestLogin(t *testing.T) {
	setup()

	responseData := postFormJSON("/login", url.Values{
		"password": {"some-password"},
	})

	if responseData["logged_in"] == true {
		t.Errorf("Expected to not be logged in yet")
	}

	postFormJSON("/create_passwords_file", url.Values{
		"password": {"some-password"},
	})

	responseData = postFormJSON("/login", url.Values{
		"password": {"some-password"},
	})
	if v, ok := responseData["logged_in"]; v == false || ok == false {
		t.Errorf("Expected to be logged in")
	}

	teardown()
}

func TestLoginWithInvalidPassword(t *testing.T) {
	setup()

	postFormJSON("/create_passwords_file", url.Values{
		"password": {"some-password"},
	})

	response := postFormJSON("/login", url.Values{
		"password": {"invalid password"},
	})

	if v, ok := response["logged_in"]; v == true || ok == false {
		t.Errorf("Expected to not be logged in")
	}

	teardown()
}

func TestLoggedIn(t *testing.T) {
	setup()

	postFormJSON("/create_passwords_file", url.Values{"password": {"some-password"}})

	response := getJSON("/logged_in")
	if v, ok := response["logged_in"]; v == true || ok == false {
		t.Errorf("Expected not to be logged in")
	}

	postFormJSON("/login", url.Values{"password": {"some-password"}})

	response = getJSON("/logged_in")
	if v, ok := response["logged_in"]; v == false || ok == false {
		t.Errorf("Expected to be logged in")
	}

	teardown()
}

func TestStoreAndSearchPassword(t *testing.T) {
	setup()

	postFormJSON("/create_passwords_file", url.Values{
		"password": {"some-password"},
	})
	postFormJSON("/login", url.Values{
		"password": {"some-password"},
	})
	postFormJSON("/store", url.Values{
		"site":     {"https://example.org/"},
		"email":    {"bla@bla.com"},
		"password": {"site password"},
	})

	response := getJSON("/search/?" + url.QueryEscape("https://example.org"))

	if response["Site"] != "https://example.org/" {
		t.Errorf("wrong site")
	}
	if response["Username"] != "bla@bla.com" {
		t.Errorf("wrong username")
	}
	if response["Password"] != "site password" {
		t.Errorf("wrong password")
	}

	teardown()
}
