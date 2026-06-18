package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

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

// requiredHeaders é o equivalente aos três @RequestHeader obrigatórios
// presentes em todos os endpoints do WsProductController Java:
// X-User-Id, X-User-Email e X-User-Type.
//
// userID e userEmail não são usados pela lógica de negócio (assim como no
// Java original, onde são recebidos mas nunca lidos no corpo dos métodos),
// mas a presença e validade deles ainda é exigida, replicando o
// comportamento de "header obrigatório" do Spring (@RequestHeader sem
// required = false faz o Spring rejeitar a requisição se o header faltar
// ou não puder ser convertido pro tipo declarado).
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
