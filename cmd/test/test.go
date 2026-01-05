// Package functionaltest is
package main

import (
	"fmt"
	"net/http"
	"os"
	"slices"

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

type debtRecord struct {
	Debtor string `json:"debtor"`
	Amount struct {
		Dollars int `json:"dollars"`
		Cents   int `json:"cents"`
	} `json:"amount"`
}

type createDebtParameters struct {
	Transactions []debtRecord `json:"transactions"`
	Creditor     string       `json:"creditor"`
	GroupID      string       `json:"group_id"`
	Title        string       `json:"title"`
}

// main is the entry point of the functional test suite
func main() {
	/*TODO: run some golden test cases here
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
				19) create user3, do some transactions between the users, make sure the record is right, settle up, make sure the record is right
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

	// 02) register user1
	user1 := createAndLoginUser(client, "test@email.com", "password", "testy mctestface")

	// 03) create group
	group1, err := createGroup(client, "group 1", user1.Token)
	if err != nil {
		log.Fatal("group creation failed", "error", err)
	}
	log.Info("created group", "group ID", group1.ID)

	// 04) join group
	err = joinGroup(client, group1.Name, user1.Token)
	if err != nil {
		log.Fatal("unable to join group")
	}
	log.Info("successfully joined the group")

	// 05) create another user and have it join the group
	user2 := createAndLoginUser(client, "test2@email.com", "password", "another clever name")

	err = joinGroup(client, group1.Name, user2.Token)
	if err != nil {
		log.Fatal("unable to join group with user 2")
	}

	log.Info("successfully created second user and joined the group")

	// 06) create another group, join with user 1, make sure user 2 does not show up as a member
	group2, err := createGroup(client, "group 2", user2.Token)
	if err != nil {
		log.Fatal("unable to create a new group")
	}

	err = joinGroup(client, group2.Name, user1.Token)
	if err != nil {
		log.Fatal("unable to join group with user 2")
	}

	users := getGroupMembers(client, group2.ID.String(), user2.Token)
	fmt.Println(users)

	// 07) create a debt for user 2
	//func createDebt(client *http.Client, groupID, token string, debtor, creditor string, title string, amount struct {
	debtID := createDebt(client, group1.ID.String(), user1.Token, user2.ID.String(), user1.ID.String(), "Debt 1", struct {
		Dollars int `json:"dollars"`
		Cents   int `json:"cents"`
	}{1, 1})
	if debtID == "" {
		log.Fatal("failed to create debt")
	}
	log.Info("successfully created debt", "debtID", debtID)

	// 07) check transaction record for group 1 and group 2
	transactionIDsg1 := getTransactions(client, user1.Token, group1.ID.String())
	log.Info("transaction in group 1", "IDs", transactionIDsg1)

	transactionIDsg2 := getTransactions(client, user1.Token, group2.ID.String())
	log.Info("transaction in group 2", "IDs", transactionIDsg2)

	if slices.Equal(transactionIDsg1, transactionIDsg2) {
		log.Fatal("group1 and group2 have the same transactions")

	}

	// 09) create a debt for user 1
	// 10) check that transaction record and balances are as expected for user 1 and user 2
	// 11) add another debt for user 1, verify everything is as expected
	// 12) delete user 1 first debt
	// 13) make sure the balance is correct
	// 14) settle up user 1
	// 15) make sure that the debt is marked as paid
	// 16) make sure that user 2s debt is not paid yet
	// 17) delete user 1s paid debt
	// 18) make sure that user 2s balance is updated
	// 19) settle up user 2, make sure that all balances are now 0
	// 20) create user3, do some transactions between the users, make sure the record is right, settle up, make sure the record is right

}
