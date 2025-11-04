// Package functionaltest is
package main

import (
	"fmt"

	"net/http"
	"os"

	"github.com/charmbracelet/log"

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
	    	01) reset db
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
	// 01) reset db
	log.Info("reset DB", "successful", reset(client, adminKey))

	//02) register user
	user1, err := register(client, registerParameters{
		Email:    "test@test.com",
		Password: "pass",
		Name:     "testy mctestface",
	})
	if err != nil {
		log.Fatal("register user failed", "error", err)
	}
	log.Info("created user", "email", user1.Email)

	user1, err = login(loginParameters{
		Email:         "test@test.com",
		Password:      "pass",
		ExpiresInSecs: 100,
	})
	if err != nil {
		log.Fatal("failed to login user 1", "error", err)
	}
}
