package dtos

import "math"

// Page é o equivalente Go do org.springframework.data.domain.Page<T>.
// Representa uma página de resultados junto com metadados de paginação,
// serializados em um formato compatível com o que o Spring Data retorna
// por padrão (campos content, totalElements, totalPages, number, size,
// first, last).
type Page[T any] struct {
	Content       []T   `json:"content"`
	TotalElements int64 `json:"totalElements"`
	TotalPages    int   `json:"totalPages"`
	Number        int   `json:"number"` // página atual (0-based), equivalente a Pageable.getPageNumber()
	Size          int   `json:"size"`
	First         bool  `json:"first"`
	Last          bool  `json:"last"`
	Empty         bool  `json:"empty"`
}

// NewPage monta um Page a partir do conteúdo da página atual e dos
// parâmetros de paginação, calculando os metadados (totalPages, first,
// last, empty) do mesmo jeito que o Spring Data faz internamente.
func NewPage[T any](content []T, page, size int, totalElements int64) *Page[T] {
	totalPages := 0
	if size > 0 {
		totalPages = int(math.Ceil(float64(totalElements) / float64(size)))
	}

	return &Page[T]{
		Content:       content,
		TotalElements: totalElements,
		TotalPages:    totalPages,
		Number:        page,
		Size:          size,
		First:         page == 0,
		Last:          page >= totalPages-1,
		Empty:         len(content) == 0,
	}
}
