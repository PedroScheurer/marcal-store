package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/PedroScheurer/product-service/internal/apperrors"
	"github.com/PedroScheurer/product-service/internal/dtos"
	"github.com/PedroScheurer/product-service/internal/services"
)

// WsProductController é o equivalente Go do WsProductController Java
// (@RestController @RequestMapping("/ws/products")).
type WsProductController struct {
	service *services.WsProductService
}

func NewWsProductController(service *services.WsProductService) *WsProductController {
	return &WsProductController{service: service}
}

// RegisterRoutes monta as rotas deste controller no router chi.
func (c *WsProductController) RegisterRoutes(r chi.Router) {
	r.Route("/ws/products", func(r chi.Router) {
		r.Post("/", c.postProduct)
		r.Put("/{idProduct}", c.putProduct)
		r.Delete("/{idProduct}", c.deleteProduct)
	})
}

// postProduct é o equivalente a WsProductController.postProduct(...).
func (c *WsProductController) postProduct(w http.ResponseWriter, r *http.Request) {
	headers, err := parseRequiredHeaders(r)
	if err != nil {
		writeWsError(w, err)
		return
	}

	var dto dtos.ProductInDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		writeWsError(w, err)
		return
	}

	productDTO, err := c.service.CreateProduct(r.Context(), dto, headers.userType)
	if err != nil {
		writeWsError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, productDTO)
}

// putProduct é o equivalente a WsProductController.putProduct(...).
func (c *WsProductController) putProduct(w http.ResponseWriter, r *http.Request) {
	idProduct, err := parseIDPathParam(r, "idProduct")
	if err != nil {
		writeWsError(w, err)
		return
	}

	headers, err := parseRequiredHeaders(r)
	if err != nil {
		writeWsError(w, err)
		return
	}

	var dto dtos.ProductUpdateInDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		writeWsError(w, err)
		return
	}

	productDTO, err := c.service.AlterProduct(r.Context(), idProduct, dto, headers.userType)
	if err != nil {
		writeWsError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, productDTO)
}

// deleteProduct é o equivalente a WsProductController.deleteProduct(...).
func (c *WsProductController) deleteProduct(w http.ResponseWriter, r *http.Request) {
	idProduct, err := parseIDPathParam(r, "idProduct")
	if err != nil {
		writeWsError(w, err)
		return
	}

	headers, err := parseRequiredHeaders(r)
	if err != nil {
		writeWsError(w, err)
		return
	}

	if err := c.service.DeleteProduct(r.Context(), idProduct, headers.userType); err != nil {
		writeWsError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// writeWsError delega ao mesmo handler global de erros usado pelo
// ProductController (apperrors.WriteErrorResponse). Isso cobre tanto os
// erros conhecidos (ProductNotFoundError -> 404, AuthenticationError -> 401)
// quanto o caso genérico, que cai no default e responde 400 com a
// mensagem "crua" do erro — equivalente ao
// @ExceptionHandler(Exception.class) handleException(...) específico
// do WsProductController Java.
func writeWsError(w http.ResponseWriter, err error) {
	apperrors.WriteErrorResponse(w, err)
}
