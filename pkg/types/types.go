package types

// Money for dirams, sents, kopeyks
type Money int64

// Category describes payment category
type PaymentCategory string

//Status info about payment status
type PaymentStatus string

//Statuses for payments
const (
	PaymentStatusOk         PaymentStatus = "OK"
	PaymentStatusFail       PaymentStatus = "FAIL"
	PaymentStatusInProgress PaymentStatus = "INPROGRESS"
)

// Payment describes structure for payments
type Payment struct {
	ID        string
	AccountID int64
	Amount    Money
	Category  PaymentCategory
	Status    PaymentStatus
}

// Phone with string type
type Phone string

// Account gives info about user's account
type Account struct {
	ID      int64
	Phone   Phone
	Balance Money
}
