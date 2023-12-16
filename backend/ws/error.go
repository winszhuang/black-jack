package ws

type WSError int

const (
	ErrForWrongRequestFormat WSError = iota
	ErrForRequestNoData
	ErrForNotCorrectState
	ErrForServerError
	ErrForDealCard
	ErrForGetAllPlayerCards
	ErrForWrongFlow
	ErrForUnauthorized
	ErrForInvalidRequest
	ErrForCantFindThisRoom
	ErrForClientNotInRoom
)
