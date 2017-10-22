package main_test

import (
	"."
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// reference to the app
var a main.App

func TestMain(m *testing.M) {
	a = main.App{}
	a.Initialize("postgres", "spoonlamp", "go-comments-api")

	ensureTableExists()

	code := m.Run()

	clearTable()

	os.Exit(code)
}

const tableCreationQuery = `
	CREATE TABLE IF NOT EXISTS comments
	(
		id SERIAL,
		username TEXT NOT NULL,
		body TEXT NOT NULL,
		mod_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (id)
	)`

func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM comments")
	a.DB.Exec("ALTER SEQUENCE comments_id_seq RESTART WITH 1")
}

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/comments/0", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestGetProduct(t *testing.T) {
	clearTable()
	addProducts(1)

	req, _ := http.NewRequest("GET", "/comment/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestCreateProduct(t *testing.T) {
	clearTable()

	payload := []byte(`{"username":"Bob", "body":"This is a comment"}`)

	req, _ := http.NewRequest("POST", "/comment", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["username"] != "Bob" {
		t.Errorf("Expected username to be 'Bob'. Got '%v'", m["username"])
	}
	if m["body"] != "This is a comment" {
		t.Errorf("Expected body to be 'This is a comment'. Got '%v'", m["body"])
	}
	if m["id"] != 1.0 {
		t.Errorf("Expected product ID to be '1'. Got '%v'", m["id"])
	}
}

func TestDeleteProduct(t *testing.T) {
	clearTable()
	addProducts(1)

	req, _ := http.NewRequest("GET", "/comment/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/comment/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/comment/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func TestUpdateProduct(t *testing.T) {
	clearTable()
	addProducts(1)

	req, _ := http.NewRequest("GET", "/comment/1", nil)
	response := executeRequest(req)
	var originalComment map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalComment)

	payload := []byte(`{"body":"updatedBody"}`)

	req, _ = http.NewRequest("PUT", "/comment/1", bytes.NewBuffer(payload))
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != originalComment["id"] {
		t.Errorf("Expected id to remain the same (%v). Got %v", originalComment["id"], m["id"])
	}
	if m["body"] == originalComment["body"] {
		t.Errorf("Expected body to change from '%v' to 'updatedBody'. Got %v", originalComment["body"], m["body"], m["body"])
	}

}

func addProducts(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		a.DB.Exec("INSERT INTO comments(username, body) VALUES ('usertest', 'bodytest')")
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
