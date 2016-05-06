package main

import (
	"github.com/gorilla/mux"
	"github.com/jameycribbs/cribbnotes/handlers/notes_handler"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestRoutes(t *testing.T) {
	data_dir := "test_data"

	r := mux.NewRouter()

	r.HandleFunc("/", makeHandler(notes_handler.Index, data_dir)).Methods("GET")
	r.HandleFunc("/new", makeHandler(notes_handler.New, data_dir)).Methods("GET")
	r.HandleFunc("/create", makeHandler(notes_handler.Create, data_dir)).Methods("POST")
	r.HandleFunc("/{id:[0-9]+}/edit", makeHandler(notes_handler.Edit, data_dir)).Methods("GET")
	r.HandleFunc("/update", makeHandler(notes_handler.Update, data_dir)).Methods("POST")
	r.HandleFunc("/{id:[0-9]+}/delete", makeHandler(notes_handler.Delete, data_dir)).Methods("GET")
	r.HandleFunc("/destroy", makeHandler(notes_handler.Destroy, data_dir)).Methods("POST")

	server := httptest.NewServer(r)
	defer server.Close()

	// Test home route
	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	}

	expected := "Test Note 1"

	actual, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(actual), expected) {
		t.Errorf("Expected the message '%s'\n", expected)
	}

	// Test new route
	resp, err = http.Get(server.URL + "/new")
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	}

	expected = "New Note"

	actual, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(actual), expected) {
		t.Errorf("Expected the message '%s'\n", expected)
	}

	// Test create route
	resp, err = http.PostForm(server.URL+"/create", url.Values{"title": {"Test Title"}, "text": {"This is a test."}})
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	}

	expected = "Test Title"

	actual, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(actual), expected) {
		t.Errorf("Expected the message '%s'\n", expected)
	}

	// Test edit route
	resp, err = http.Get(server.URL + "/1/edit")
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	}

	expected = "Editing Note"

	actual, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(actual), expected) {
		t.Errorf("Expected the message '%s'\n", expected)
	}

	// Test update route
	resp, err = http.PostForm(server.URL+"/update", url.Values{"fileId": {"2"}, "title": {"Updated Test Title"},
		"text": {"This is a test."}})
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	}

	expected = "Updated Test Title"

	actual, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(actual), expected) {
		t.Errorf("Expected the message '%s'\n", expected)
	}

	// Test delete route
	resp, err = http.Get(server.URL + "/2/delete")
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	}

	expected = "Deleting Note"

	actual, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(actual), expected) {
		t.Errorf("Expected the message '%s'\n", expected)
	}

	// Test destroy route
	resp, err = http.PostForm(server.URL+"/destroy", url.Values{"fileId": {"2"}})
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	}
}
