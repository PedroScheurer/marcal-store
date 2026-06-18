package dtos

// ErrorResponse é o equivalente Go do record Java ErrorResponse,
// usado pelo handler global de erros para padronizar respostas de erro.
type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
