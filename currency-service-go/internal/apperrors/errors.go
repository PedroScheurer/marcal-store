package apperrors

type CurrencyNotFoundError struct {
	Message string
}

func NewCurrencyNotFoundError(message string) *CurrencyNotFoundError {
	return &CurrencyNotFoundError{Message: message}
}

func (e *CurrencyNotFoundError) Error() string {
	return e.Message
}
