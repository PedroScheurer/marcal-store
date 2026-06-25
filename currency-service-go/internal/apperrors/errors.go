package apperrors

// ProductNotFoundError é o equivalente Go da ProductNotFoundException.
// Implementa a interface error.
type ProductNotFoundError struct {
	Message string
}

func NewProductNotFoundError(message string) *ProductNotFoundError {
	return &ProductNotFoundError{Message: message}
}

func (e *ProductNotFoundError) Error() string {
	return e.Message
}

// ExternalServiceError é o equivalente Go da ExternalServiceException,
// usada quando uma dependência externa (ex.: currency-service) falha.
type ExternalServiceError struct {
	Message string
}

func NewExternalServiceError(message string) *ExternalServiceError {
	return &ExternalServiceError{Message: message}
}

func (e *ExternalServiceError) Error() string {
	return e.Message
}

// AuthenticationError é o equivalente Go da javax.naming.AuthenticationException
// usada no Java original para indicar falta de permissão do usuário.
type AuthenticationError struct {
	Message string
}

func NewAuthenticationError(message string) *AuthenticationError {
	return &AuthenticationError{Message: message}
}

func (e *AuthenticationError) Error() string {
	return e.Message
}
