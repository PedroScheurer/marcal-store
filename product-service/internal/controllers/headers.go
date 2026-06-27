package controllers

import (
	"net/http"
	"strconv"

	"github.com/PedroScheurer/product-service/internal/apperrors"
)

const adminType = 0

type requiredHeaders struct {
	userID    int64
	userEmail string
	userType  int
}

func parseRequiredHeaders(r *http.Request) (requiredHeaders, error) {
	userIDRaw := r.Header.Get("X-User-Id")
	userEmail := r.Header.Get("X-User-Email")
	userTypeRaw := r.Header.Get("X-User-Type")

	if userIDRaw == "" || userEmail == "" || userTypeRaw == "" {
		return requiredHeaders{}, apperrors.NewAuthenticationError(
			"Headers obrigatórios ausentes: X-User-Id, X-User-Email, X-User-Type",
		)
	}

	userID, err := strconv.ParseInt(userIDRaw, 10, 64)
	if err != nil {
		return requiredHeaders{}, apperrors.NewAuthenticationError("X-User-Id inválido")
	}

	userType, err := strconv.Atoi(userTypeRaw)
	if err != nil {
		return requiredHeaders{}, apperrors.NewAuthenticationError("X-User-Type inválido")
	}

	return requiredHeaders{userID: userID, userEmail: userEmail, userType: userType}, nil
}
