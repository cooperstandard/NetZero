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

	if health(client) != 204 {
		return
	}

	// 01) reset db
	log.Info("reset DB", "successful", reset(client, adminKey))

	// 02) register user
	user := createAndLoginUser(client, "test@email.com", "password", "testy mctestface")

	// 03) create group
	group, err := createGroup(client, "group 1", user.Token)
	if err != nil {
		log.Fatal("group creation failed", "error", err)
	}
	log.Info("created group", "group ID", group.ID)

	// 04) join group
	err = joinGroup(client, "group 1", user.Token)

	if err != nil {
		log.Fatal("unable to join group")
	}
	log.Info("successfully joined the group")
}

