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

func TestFindPaymentByID_fail(t *testing.T) {
	svc := &Service{}

	svc.RegisterAccount("+992900010203")
	svc.Deposit(1, 500)
	svc.Pay(1, 100, "food")

	_, err := svc.FindPaymentByID("120")

	if err != nil {
		t.Error(err)
		return
	}
}

func TestFindPaymentByID_success(t *testing.T) {
	svc := &Service{}

	svc.RegisterAccount("+992900010203")
	svc.Deposit(1, 500)
	svc.Pay(1, 100, "food")

	paymentID := svc.payments[0].ID

	payment, err := svc.FindPaymentByID(paymentID)

	if err != nil {
		t.Error(err)
		return
	}

	expected := &types.Payment{
		ID:        paymentID,
		AccountID: 1,
		Amount:    100,
		Category:  "food",
		Status:    types.PaymentStatusInProgress,
	}

	if !reflect.DeepEqual(expected, payment) {
		t.Errorf("expected %v: got this %v", expected, payment)
	}
}

func TestReject_fail(t *testing.T) {
	svc := &Service{}

	svc.RegisterAccount("+992900010203")
	svc.Deposit(1, 500)
	svc.Pay(1, 100, "food")

	err := svc.Reject("12")

	if err != nil {
		t.Error(err)
		return
	}
}

func TestReject_success(t *testing.T) {
	svc := &Service{}

	svc.RegisterAccount("+992900010203")
	svc.Deposit(1, 500)
	svc.Pay(1, 100, "food")

	paymentID := svc.payments[0].ID

	err := svc.Reject(paymentID)

	if err != nil {
		t.Error(err)
		return
	}

	expected := &types.Payment{
		ID:        paymentID,
		AccountID: 1,
		Amount:    100,
		Category:  "food",
		Status:    types.PaymentStatusFail,
	}

	payment := svc.payments[0]

	if !reflect.DeepEqual(expected, payment) {
		t.Errorf("expected %v: got this %v", expected, payment)
	}
}
