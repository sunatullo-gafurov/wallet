package wallet

import (
	"reflect"
	"testing"

	"github.com/sunatullo-gafurov/wallet/pkg/types"
)

func TestService_FindAccountByID_success(t *testing.T) {
	s := newTService()

	account, _, err := s.addAccount(defaultTAccount)
	if err != nil {
		t.Error(err)
		return
	}

	got, err := s.FindAccountByID(account.ID)
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(account, got) {
		t.Errorf("expected %v: got this %v", account, got)
		return
	}
}

func TestService_FindAccountByID_fail(t *testing.T) {
	s := newTService()

	_, _, err := s.addAccount(defaultTAccount)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = s.FindAccountByID(1484)
	if err == nil {
		t.Errorf("FindAccountByID() must return error, returned nil")
		return
	}

	if err != ErrAccountNotFound {
		t.Errorf("FindAccountByID() must return ErrAccountNotFound, returned = %v", err)
		return
	}
}

func TestService_FindPaymentByID_success(t *testing.T) {
	s := newTService()

	_, payments, err := s.addAccount(defaultTAccount)
	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]

	got, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("FindPaymentByID(): error = %v", err)
		return
	}

	if !reflect.DeepEqual(payment, got) {
		t.Errorf("FindPaymentByID(): expected %v, got this %v", payment, got)
	}
}

func TestService_FindPaymentByID_fail(t *testing.T) {
	s := newTService()

	_, _, err := s.addAccount(defaultTAccount)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = s.FindPaymentByID("xyz")
	if err == nil {
		t.Errorf("FindPaymentByID() must return error, returned  nil")
		return
	}

	if err != ErrPaymentNotFound {
		t.Errorf("FindPaymentByID() must return ErrPaymentNotFound, returned = %v", err)
	}
}

func TestService_Reject_success(t *testing.T) {
	s := newTService()

	_, payments, err := s.addAccount(defaultTAccount)
	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]
	err = s.Reject(payment.ID)
	if err != nil {
		t.Errorf("Reject(): error = %v", err)
		return
	}

	savedPayment, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("Reject(): can't find payment by id, error = %v", err)
		return
	}
	if savedPayment.Status != types.PaymentStatusFail {
		t.Errorf("Reject(): status didn't change, payment = %v", savedPayment)
		return
	}

	savedAccount, err := s.FindAccountByID(payment.AccountID)
	if err != nil {
		t.Errorf("Reject(): can't find account by id, error = %v", err)
		return
	}
	if savedAccount.Balance != defaultTAccount.balance {
		t.Errorf("Reject(): balance didn't change, account = %v", savedAccount)
		return
	}
}

func TestService_Reject_fail(t *testing.T) {
	s := newTService()

	_, payments, err := s.addAccount(defaultTAccount)
	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]
	err = s.Reject("xyz")
	if err == nil {
		t.Errorf("Reject(): must give error returned nil")
		return
	}

	savedPayment, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("Reject(): can't find payment by id, error = %v", err)
		return
	}
	if savedPayment.Status != types.PaymentStatusInProgress {
		t.Errorf("Reject(): status should not change")
		return
	}

	_, err = s.FindAccountByID(payment.AccountID)
	if err != nil {
		t.Errorf("Reject(): can't find account by id, error = %v", err)
		return
	}
}

func TestService_Repeat_success(t *testing.T) {
	s := newTService()

	_, payments, err := s.addAccount(defaultTAccount)
	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]

	got, repeatErr := s.Repeat(payment.ID)
	if repeatErr != nil {
		t.Error(repeatErr)
		return
	}

	if payment.AccountID != got.AccountID &&
		payment.Amount != got.Amount &&
		payment.Category != got.Category &&
		payment.Status != got.Status {
		t.Errorf("%v, %v", payment, got)
	}
}

func TestService_Repeat_fail(t *testing.T) {
	s := newTService()

	_, _, err := s.addAccount(defaultTAccount)
	if err != nil {
		t.Error(err)
		return
	}

	_, repeatErr := s.Repeat("xyz")
	if repeatErr == nil {
		t.Errorf("Repeat() returned: %v, should return nil", repeatErr)
		return
	}
}

func TestService_FavoritePayment_success(t *testing.T) {
	s := newTService()

	_, payments, err := s.addAccount(defaultTAccount)
	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]

	favoritePayment, err := s.FavoritePayment(payment.ID, "me")
	if err != nil {
		t.Error(err)
		return
	}

	got, err := s.FindFavoriteByID(favoritePayment.ID)
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(favoritePayment, got) {
		t.Errorf("expected %v, got this %v", favoritePayment, got)
		return
	}
}

func TestService_FavoritePayment_fail(t *testing.T) {
	s := newTService()

	_, _, err := s.addAccount(defaultTAccount)
	if err != nil {
		t.Error(err)
		return
	}

	_, favErr := s.FavoritePayment("xyz", "me")
	if favErr == nil {
		t.Errorf("FavoritePayment() should return nil")
		return
	}
}

func TestService_PayFromFavorite_success(t *testing.T) {
	s := newTService()

	_, payments, err := s.addAccount(defaultTAccount)
	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]

	favorite, err := s.FavoritePayment(payment.ID, "me")
	if err != nil {
		t.Error(err)
		return
	}

	paymentFromFavorite, err := s.PayFromFavorite(favorite.ID)
	if err != nil {
		t.Error(err)
		return
	}

	got, err := s.FindPaymentByID(paymentFromFavorite.ID)
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(paymentFromFavorite, got) {
		t.Errorf("expected %v, got this %v", paymentFromFavorite, got)
		return
	}
}

func TestService_PayFromFavorite_fail(t *testing.T) {
	s := newTService()

	_, payments, err := s.addAccount(defaultTAccount)
	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]

	_, favErr := s.FavoritePayment(payment.ID, "me")
	if err != nil {
		t.Error(favErr)
		return
	}

	_, payErr := s.PayFromFavorite("xyz")
	if payErr == nil {
		t.Errorf("PayFromFavorite() should return nil")
		return
	}
}
