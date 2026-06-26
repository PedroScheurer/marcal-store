package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/PedroScheurer/currency-service/internal/apperrors"
	"github.com/PedroScheurer/currency-service/internal/services"
)

// CurrencyController é o equivalente Go do CurrencyController Java
// (@RestController @RequestMapping("/currency")).
type CurrencyController struct {
	service *services.CurrencyService
}

func NewCurrencyController(service *services.CurrencyService) *CurrencyController {
	return &CurrencyController{service: service}
}

func (c *CurrencyController) RegisterRoutes(r chi.Router) {
	r.Route("/currency", func(r chi.Router) {
		// Equivalente a @GetMapping(path = "/convert", params = {"source", "target"})
		r.Get("/convert", c.convert)
	})
}

// convert é o equivalente a CurrencyController.findBySourceAndTarget.
func (c *CurrencyController) convert(w http.ResponseWriter, r *http.Request) {
	source := r.URL.Query().Get("source")
	target := r.URL.Query().Get("target")

	if source == "" || target == "" {
		http.NotFound(w, r)
		return
	}

	dto, err := c.service.FindBySourceAndTarget(r.Context(), source, target)
	if err != nil {
		apperrors.WriteErrorResponse(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(dto)
}
