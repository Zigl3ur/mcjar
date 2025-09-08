package utils

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetReqJsonData(t *testing.T) {
	var jsonData = map[string]string{
		"message": "success",
		"data":    "test",
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		jsonBytes, _ := json.Marshal(jsonData)
		w.Write(jsonBytes)
	}))

	defer mockServer.Close()

	var respData map[string]string
	_, err := GetReqJson(mockServer.URL, &respData)
	if err != nil {
		t.Fatalf("got error %v, expected no error", err)
	}
	if respData["message"] != "success" || respData["data"] != "test" {
		t.Errorf("got JSON data: %v, expected {\"message\": \"success\", \"data\": \"test\"}", respData)
	}

}

func TestGetReqJsonOK(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{}`))
	}))

	defer mockServer.Close()

	var respData map[string]string
	status, err := GetReqJson(mockServer.URL, &respData)
	if err != nil {
		t.Fatalf("got error %v, expected no error", err)
	}
	if status != http.StatusOK {
		t.Errorf("got status %v, expected status 200", status)
	}
}

func TestGetReqJsonError(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))

	defer mockServer.Close()

	var respData map[string]string
	status, err := GetReqJson(mockServer.URL, &respData)
	if err == nil {
		t.Fatalf("got no error, expected error")
	}
	if status != http.StatusInternalServerError {
		t.Errorf("got status %v, expected status 500", status)
	}
}

func TestGetReqXmlData(t *testing.T) {
	type TestData struct {
		Message string `xml:"message"`
		Data    string `xml:"data"`
	}

	xmlData := TestData{
		Message: "success",
		Data:    "test",
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		xmlBytes, _ := xml.Marshal(xmlData)
		w.Write(xmlBytes)
	}))

	defer mockServer.Close()

	var respData TestData
	_, err := GetReqXml(mockServer.URL, &respData)
	if err != nil {
		t.Fatalf("got %v, expected no error", err)
	}
	if respData.Message != "success" || respData.Data != "test" {
		t.Errorf("got XML data: %v, expected {\"message\": \"success\", \"data\": \"test\"}", respData)
	}
}

func TestGetReqXmlOK(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/xml")
		w.Write([]byte(`<root></root>`))
	}))

	defer mockServer.Close()

	type EmptyRoot struct{}
	var respData EmptyRoot
	status, err := GetReqXml(mockServer.URL, &respData)
	if err != nil {
		t.Fatalf("got %v, expected no error", err)
	}
	if status != http.StatusOK {
		t.Errorf("got status %v, expected status 200", status)
	}
}

func TestGetReqXmlError(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))

	defer mockServer.Close()

	type TestData struct {
		Message string `xml:"message"`
	}
	var respData TestData
	status, err := GetReqXml(mockServer.URL, &respData)
	if err == nil {
		t.Fatalf("got no error, expected an error")
	}
	if status != http.StatusInternalServerError {
		t.Errorf("got status %v, expected status 500", status)
	}
}
