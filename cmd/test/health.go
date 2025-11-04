package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/cooperstandard/NetZero/internal/routes"
	"io"
	"net/http"
)

func health(client *http.Client) int {
	_, status := doRequest(client, "GET", "/health", nil, "")
	return status
}

func register(client *http.Client, params registerParameters) (routes.User, error) {
	body, _ := json.Marshal(params)
	resp, status := doRequest(client, "POST", "/register", body, "")
	if status != 201 && status != 200 {
		log.Error("status for request", "status", status)
		return routes.User{}, fmt.Errorf("unable to register user: %s", params.Email)
	}

	var user routes.User

	respBody, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(respBody, &user)

	return user, nil
}

func login(params loginParameters) (routes.User, error) {
	return routes.User{}, nil
}

func reset(client *http.Client, key string) bool {
	_, status := doRequest(client, "POST", "/admin/reset", nil, key)
	return status != 0 && status < 300
}

func doRequest(client *http.Client, method string, endpoint string, body []byte, token string) (*http.Response, int) {
	var req *http.Request
	var err error
	if len(body) != 0 {
		req, err = http.NewRequest(method, basepath+endpoint, bytes.NewBuffer(body))
	} else {
		req, err = http.NewRequest(method, basepath+endpoint, nil)
	}
	if err != nil {
		return nil, 0
	}
	if token != "" {
		req.Header.Add("Authorization", "Bearer: "+token)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, 0
	}

	return res, res.StatusCode
}
