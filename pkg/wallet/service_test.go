package wallet

import (
	"reflect"
	"testing"

	"github.com/sunatullo-gafurov/wallet/pkg/types"
)

func TestFindAccountByID(t *testing.T) {
	svc := &Service{}

	svc.RegisterAccount("+992900010203")
	svc.RegisterAccount("+992900040203")
	svc.RegisterAccount("+992900050203")

	acc, err := svc.FindAccountByID(1)

	if err != nil {
		t.Error(err)
		return
	}

	expected := &types.Account{
		ID:      1,
		Phone:   "+992900010203",
		Balance: 0,
	}

	if !reflect.DeepEqual(expected, acc) {
		t.Errorf("expected %v: got this %v", expected, acc)
	}
}
