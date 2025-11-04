// Package functionaltest is
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/charmbracelet/log"
	"github.com/cooperstandard/NetZero/internal/routes"
	"github.com/joho/godotenv"
)

const basepath string = "http://localhost:8080/api/v1"

type registerParameters struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type loginParameters struct {
	Email         string `json:"email"`
	Password      string `json:"password"`
	ExpiresInSecs int    `json:"expiresInSeconds"` // TODO: for testing, configure this in the environment
}

// main is the entry point of the functional test suite
func main() {
	/*TODO: run some golden test cases here
				00) health
	    	01) reset db
	  		02) register user
	    	2a) login
	    	03) create group
	  		04) join group
	  		05) create another user and have it join the group
	  		06) create another group, join with user 1, make sure user 2 does not show up as a member
	    	06) create a debt for user 2
	  		07) check transaction record for user 1 and user 2
	  		08) create a debt for user 1
	    	09) check that transaction record and balances are as expected for user 1 and user 2
	  		10) add another debt for user 1, verify everything is as expected
	  		11) delete user 1 first debt
	  		12) make sure the balance is correct
	  		13) settle up user 1
	  		14) make sure that the debt is marked as paid
	  		15) make sure that user 2s debt is not paid yet
	  		16) delete user 1s paid debt
	  		17) make sure that user 2s balance is updated
	  		18) settle up user 2, make sure that all balances are now 0
	*/

	log.Warn("starting functional tests")

	godotenv.Load()

	adminKey := os.Getenv("ADMIN_KEY")

	client := &http.Client{}
	log.Info(fmt.Sprintf("received code %d from call to health", health(client)))
	log.Info("reset DB", "successful", reset(client, adminKey))
	user1, err := register(client, registerParameters{
		Email:    "test@test.com",
		Password: "pass",
		Name:     "testy mctestface",
	})
	if err != nil {
		log.Fatal("register user failed", "error", err)
	}
	log.Info("created user", "email", user1.Email)
}

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
