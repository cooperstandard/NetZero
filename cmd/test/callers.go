package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/cooperstandard/NetZero/internal/routes"
)

func health(client *http.Client) int {
	_, status := doRequest(client, "GET", "/health", nil, "")
	return status
}

func createGroup(client *http.Client, groupName string, token string) (routes.Group, error) {
	params, err := json.Marshal(struct {
		Name string `json:"name"`
	}{Name: groupName})

	if err != nil {
		return routes.Group{}, err
	}

	response, status := doRequest(client, "POST", "/groups", params, token)
	if status != 200 {
		fmt.Printf("recieved status %d while trying to create group\n", status)
		return routes.Group{}, fmt.Errorf("unable to create group with name %s", groupName)
	}

	var group routes.Group

	respBody, _ := io.ReadAll(response.Body)
	json.Unmarshal(respBody, &group)

	return group, nil
}

func createDebt(client *http.Client, groupID, token string) string {



	return ""
}

func getGroupMembers(client *http.Client, groupID, token string) []string {
	resp, status := doRequest(client, "GET", "/groups/members/" + groupID, nil, token)
	if status != 200 {
		return nil
	}

	var users []routes.User
	respBody, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(respBody, &users)

	var ids []string
	for _, v := range users {
		ids = append(ids, v.ID.String())
	}


	return ids
}

func joinGroup(client *http.Client, groupName string, token string) error {
	params, err := json.Marshal(struct {
		Name string `json:"group_name"`
	}{Name: groupName})

	if err != nil {
		return err
	}

	_, status := doRequest(client, "POST", "/groups/join", params, token)
	
	if status != 204 {
		return fmt.Errorf("unable to join group")
	}

	return nil
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

func login(client *http.Client, params loginParameters) (routes.User, int) {
	body, _ := json.Marshal(params)
	resp, status := doRequest(client, "POST", "/login", body, "")
	if status != 200 {
		log.Error("login failed with", "status", status)
		return routes.User{}, status
	}

	var user routes.User

	respBody, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(respBody, &user)

	return user, status
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

