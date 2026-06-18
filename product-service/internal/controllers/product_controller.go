package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/PedroScheurer/product-service/internal/apperrors"
	"github.com/PedroScheurer/product-service/internal/services"
)

// ProductController é o equivalente Go do ProductController Java
// (@RestController @RequestMapping("/products")).
type ProductController struct {
	service *services.ProductService
}

func NewProductController(service *services.ProductService) *ProductController {
	return &ProductController{service: service}
}

// RegisterRoutes monta as rotas deste controller no router chi,
// equivalente às anotações @GetMapping da classe Java.
func (c *ProductController) RegisterRoutes(r chi.Router) {
	r.Route("/products", func(r chi.Router) {
		// Equivalente a:
		// @GetMapping(path = {"{idProduct}"}, params = {"targetCurrency"})
		// Como o chi não tem o conceito de "params obrigatórios" do Spring
		// MVC, fazemos essa checagem manualmente dentro do handler.
		r.Get("/{idProduct}", c.findProduct)

		// @GetMapping("/noconverter/{idProduct}")
		r.Get("/noconverter/{idProduct}", c.getProductNoConverter)

		// @GetMapping (lista paginada)
		r.Get("/", c.getAllProducts)
	})
}

// findProduct é o equivalente a ProductController.findProduct(idProduct, targetCurrency).
func (c *ProductController) findProduct(w http.ResponseWriter, r *http.Request) {
	idProduct, err := parseIDPathParam(r, "idProduct")
	if err != nil {
		apperrors.WriteErrorResponse(w, err)
		return
	}

	targetCurrency := r.URL.Query().Get("targetCurrency")
	if targetCurrency == "" {
		// Replica o comportamento do `params = {"targetCurrency"}` do Spring:
		// sem esse parâmetro, a rota nem deveria casar. Aqui retornamos 404
		// como o Spring faz quando nenhum handler casa com a requisição.
		http.NotFound(w, r)
		return
	}

	productDTO, err := c.service.FindByID(r.Context(), idProduct, targetCurrency)
	if err != nil {
		apperrors.WriteErrorResponse(w, err)
		return
	}

	writeJSON(w, http.StatusOK, productDTO)
}

// getProductNoConverter é o equivalente a ProductController.getProductNoConverter(idProduct).
func (c *ProductController) getProductNoConverter(w http.ResponseWriter, r *http.Request) {
	idProduct, err := parseIDPathParam(r, "idProduct")
	if err != nil {
		apperrors.WriteErrorResponse(w, err)
		return
	}

	dto, err := c.service.FindProductNoConversion(r.Context(), idProduct)
	if err != nil {
		apperrors.WriteErrorResponse(w, err)
		return
	}

	writeJSON(w, http.StatusOK, dto)
}

// getAllProducts é o equivalente a ProductController.getAllProducts(targetCurrency, pageable).
func (c *ProductController) getAllProducts(w http.ResponseWriter, r *http.Request) {
	targetCurrency := r.URL.Query().Get("targetCurrency")
	if targetCurrency == "" {
		http.NotFound(w, r)
		return
	}

	page, size, sortBy, sortDir := parsePageable(r)

	productDTOs, err := c.service.FindProductsPaged(r.Context(), targetCurrency, page, size, sortBy, sortDir)
	if err != nil {
		apperrors.WriteErrorResponse(w, err)
		return
	}

	writeJSON(w, http.StatusOK, productDTOs)
}

// parsePageable lê os parâmetros de paginação da query string, replicando
// o @PageableDefault(page = 0, size = 5, sort = "description", direction = ASC)
// do controller Java. O Spring aceita ?page=&size=&sort=campo,direção.
func parsePageable(r *http.Request) (page, size int, sortBy, sortDir string) {
	page = 0
	size = 5
	sortBy = "description"
	sortDir = "ASC"

	q := r.URL.Query()

	if v := q.Get("page"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed >= 0 {
			page = parsed
		}
	}

	if v := q.Get("size"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
			size = parsed
		}
	}

	if v := q.Get("sort"); v != "" {
		// Spring aceita "campo,direção" (ex.: "price,desc") no mesmo parâmetro.
		field, dir := splitSortParam(v)
		sortBy = field
		if dir != "" {
			sortDir = dir
		}
	}

	return page, size, sortBy, sortDir
}

func splitSortParam(v string) (field, dir string) {
	for i := 0; i < len(v); i++ {
		if v[i] == ',' {
			return v[:i], v[i+1:]
		}
	}
	return v, ""
}

func parseIDPathParam(r *http.Request, name string) (int64, error) {
	raw := chi.URLParam(r, name)
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return 0, apperrors.NewProductNotFoundError("ID de produto inválido: " + raw)
	}
	return id, nil
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
