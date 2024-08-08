package api

import (
	"fmt"
	"reflect"
	"testing"

	db "github.com/arkarsg/splitapp/db/sqlc"
	u "github.com/arkarsg/splitapp/utils"
)

func TestCreateAccount(t *testing.T) {

}

// Custom gomock matcher
type AccountMatcher struct {
	args          db.CreateAccountParams
	plainPassword string
}

func (a *AccountMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateAccountParams)
	if !ok {
		return false
	}

	ok = u.CheckPasswordHash(a.plainPassword, arg.HashedPassword)
	if !ok {
		return false
	}

	a.args.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(a.args, arg)
}

func (a *AccountMatcher) String() string {
	return fmt.Sprintf("Matches %v with %v", a.args, a.plainPassword)
}
