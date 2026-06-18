package apperrors

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/PedroScheurer/product-service/internal/dtos"
)

// WriteErrorResponse é o equivalente Go do GlobalExceptionHandler (@ControllerAdvice).
//
// Em Spring, o ControllerAdvice intercepta exceções lançadas por qualquer
// controller/service e converte para uma resposta HTTP padronizada. Go não
// tem esse mecanismo de interceptação automática, então cada handler chama
// esta função explicitamente no seu bloco de tratamento de erro:
//
//	if err != nil {
//	    apperrors.WriteErrorResponse(w, err)
//	    return
//	}
func WriteErrorResponse(w http.ResponseWriter, err error) {
	var notFound *ProductNotFoundError
	var external *ExternalServiceError
	var auth *AuthenticationError

	switch {
	case errors.As(err, &notFound):
		writeJSON(w, http.StatusNotFound, dtos.ErrorResponse{
			Status:  http.StatusNotFound,
			Message: notFound.Message,
		})
	case errors.As(err, &external):
		writeJSON(w, http.StatusNotFound, dtos.ErrorResponse{
			Status:  http.StatusNotFound,
			Message: external.Message,
		})
	case errors.As(err, &auth):
		writeJSON(w, http.StatusUnauthorized, dtos.ErrorResponse{
			Status:  http.StatusUnauthorized,
			Message: auth.Message,
		})
	default:
		// Equivalente ao handler genérico Exception.class do WsProductController:
		// retorna 400 com a mensagem de erro "crua".
		writeJSON(w, http.StatusBadRequest, dtos.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
