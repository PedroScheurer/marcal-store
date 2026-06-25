package infra

import (
	"encoding/json"
	"net/http"
)

// AppInfo representa os dados que o endpoint /info vai expor
type AppInfo struct {
	App struct {
		Name        string `json:"name"`
		Version     string `json:"version"`
		Description string `json:"description"`
		Author      string `json:"author"`
	} `json:"app"`
}

// NewInfoHandler cria um handler HTTP que responde com os dados da aplicação
func NewInfoHandler() http.HandlerFunc {
	info := AppInfo{}
	info.App.Name = "product-service"
	info.App.Version = "2.0.0"
	info.App.Description = "Microsserviço de produtos migrado de Java para Go"
	info.App.Author = "Pedro Konig Scheurer"

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(info)
	}
}
