// Package functionaltest is
package functionaltest

import "time"
import "github.com/charmbracelet/log"

//Start is the entry point of the functional test suite
func Start() bool{
/*TODO: run some golden test cases here
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
		19) 
*/
	log.Info("test")

	//let the server start
	time.Sleep(10 * time.Second)

	//stop the server
	panic(1)

}

