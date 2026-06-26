package apperrors

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/PedroScheurer/currency-service/internal/dtos"
)

func WriteErrorResponse(w http.ResponseWriter, err error) {
	var notFound *CurrencyNotFoundError

	switch {
	case errors.As(err, &notFound):
		writeJSON(w, http.StatusNotFound, dtos.ErrorResponse{
			Status:  http.StatusNotFound,
			Message: notFound.Message,
		})
	default:
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
