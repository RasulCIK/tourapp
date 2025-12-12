package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

const baseURL = "http://localhost:8080"

var jwtToken string


func TestRegisterUser(t *testing.T) {
	url := fmt.Sprintf("%s/register", baseURL)
	payload := map[string]string{
		"username": "rassul",
		"email":    "rassul@gmail.com",
		"password": "123456",
	}
	data, _ := json.Marshal(payload)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(resp.Body)
		t.Fatalf("Expected 201, got %d: %s", resp.StatusCode, string(body))
	}
}


func TestLoginUser(t *testing.T) {
	url := fmt.Sprintf("%s/login", baseURL)
	payload := map[string]string{
		"email":    "rassul@gmail.com",
		"password": "123456",
	}
	data, _ := json.Marshal(payload)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)

	token, ok := result["token"]
	if !ok {
		t.Fatal("Token not found in response")
	}

	jwtToken = token 
}


func TestGetUser(t *testing.T) {
	if jwtToken == "" {
		t.Fatal("JWT token is empty")
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/users/1", baseURL), nil)
	req.Header.Set("Authorization", "Bearer "+jwtToken)

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		t.Fatalf("Expected 200, got %d: %s", resp.StatusCode, string(body))
	}
}


func TestUpdateUser(t *testing.T) {
	client := &http.Client{}
	payload := map[string]string{
		"username": "rassul_updated",
		"email":    "rassul_updated@example.com",
		"password": "newpassword123",
	}
	data, _ := json.Marshal(payload)

	req, _ := http.NewRequest("PUT", fmt.Sprintf("%s/users/1", baseURL), bytes.NewBuffer(data))
	req.Header.Set("Authorization", "Bearer "+jwtToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		t.Fatalf("Expected 200, got %d: %s", resp.StatusCode, string(body))
	}
}


func TestDeleteUser(t *testing.T) {
	client := &http.Client{}
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/users/1", baseURL), nil)
	req.Header.Set("Authorization", "Bearer "+jwtToken)

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		t.Fatalf("Expected 200, got %d: %s", resp.StatusCode, string(body))
	}
}
