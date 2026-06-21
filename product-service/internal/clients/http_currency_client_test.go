package clients

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"
)

// newFakeCurrencyService sobe um httptest.Server simulando o
// currency-service real, respondendo GET /currency/convert.
// handler permite customizar o comportamento por teste (status, body, atraso).
func newFakeCurrencyService(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	return server
}

// newFakeEureka sobe um httptest.Server simulando o Eureka Server,
// respondendo GET /apps/{APP_NAME} com uma única instância UP apontando
// para o endereço do currencyServer informado.
func newFakeEureka(t *testing.T, currencyServer *httptest.Server) *httptest.Server {
	t.Helper()

	targetURL, err := url.Parse(currencyServer.URL)
	if err != nil {
		t.Fatalf("failed to parse currency server url: %v", err)
	}

	host, portStr, err := net.SplitHostPort(targetURL.Host)
	if err != nil {
		t.Fatalf("failed to split host/port: %v", err)
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		t.Fatalf("failed to parse port: %v", err)
	}

	eureka := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(EurekaAppResponse{
			Application: EurekaApplication{
				Name: "CURRENCY-SERVICE",
				Instances: []EurekaInstanceQuery{
					{
						HostName: host,
						Status:   "UP",
						Port:     EurekaPortQuery{Number: port, Enabled: "true"},
					},
				},
			},
		})
	}))
	t.Cleanup(eureka.Close)
	return eureka
}

func TestHTTPCurrencyClient_GetCurrency_Success(t *testing.T) {
	currencyServer := newFakeCurrencyService(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/currency/convert" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("source") != "BRL" || r.URL.Query().Get("target") != "USD" {
			t.Fatalf("unexpected query params: %s", r.URL.RawQuery)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(CurrencyResponse{
			SourceCurrency: "BRL",
			TargetCurrency: "USD",
			ConversionRate: 0.2,
			Environment:    "currency-service test",
		})
	})

	eureka := newFakeEureka(t, currencyServer)

	discovery := NewServiceDiscovery(eureka.URL, &http.Client{Timeout: 2 * time.Second}, 10*time.Second)
	client := NewHTTPCurrencyClient(discovery, 2*time.Second)

	result, err := client.GetCurrency(context.Background(), "BRL", "USD")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ConversionRate != 0.2 {
		t.Fatalf("expected rate 0.2, got %v", result.ConversionRate)
	}
}

func TestHTTPCurrencyClient_GetCurrency_NotFound(t *testing.T) {
	attempts := 0
	currencyServer := newFakeCurrencyService(t, func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusNotFound)
	})
	eureka := newFakeEureka(t, currencyServer)

	discovery := NewServiceDiscovery(eureka.URL, &http.Client{Timeout: 2 * time.Second}, 10*time.Second)
	client := NewHTTPCurrencyClient(discovery, 2*time.Second)

	_, err := client.GetCurrency(context.Background(), "BRL", "XYZ")

	if err == nil {
		t.Fatal("expected error for 404 response, got nil")
	}
	if !errors.Is(err, ErrCurrencyNotFound) {
		t.Fatalf("expected ErrCurrencyNotFound, got: %v", err)
	}

	// Erro de negócio não deve disparar retry: só 1 tentativa esperada.
	if attempts != 1 {
		t.Fatalf("expected exactly 1 attempt for business error (no retry), got %d", attempts)
	}
}

func TestHTTPCurrencyClient_GetCurrency_RetriesOnServerError(t *testing.T) {
	attempts := 0

	currencyServer := newFakeCurrencyService(t, func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(CurrencyResponse{
			SourceCurrency: "BRL",
			TargetCurrency: "USD",
			ConversionRate: 0.2,
		})
	})
	eureka := newFakeEureka(t, currencyServer)

	discovery := NewServiceDiscovery(eureka.URL, &http.Client{Timeout: 2 * time.Second}, 10*time.Second)
	client := NewHTTPCurrencyClient(discovery, 2*time.Second)

	result, err := client.GetCurrency(context.Background(), "BRL", "USD")
	if err != nil {
		t.Fatalf("unexpected error after retries: %v", err)
	}
	if result.ConversionRate != 0.2 {
		t.Fatalf("expected rate 0.2, got %v", result.ConversionRate)
	}
	if attempts != 3 {
		t.Fatalf("expected 3 attempts, got %d", attempts)
	}
}

func TestServiceDiscovery_ResolveURL_UsesStaleCacheWhenEurekaDown(t *testing.T) {
	currencyServer := newFakeCurrencyService(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(CurrencyResponse{ConversionRate: 0.2})
	})
	eureka := newFakeEureka(t, currencyServer)

	discovery := NewServiceDiscovery(eureka.URL, &http.Client{Timeout: 2 * time.Second}, 10*time.Second)

	// Primeira resolução: Eureka está no ar, preenche o cache.
	firstURL, err := discovery.ResolveURL(context.Background(), "currency-service")
	if err != nil {
		t.Fatalf("unexpected error on first resolve: %v", err)
	}

	// Derruba o Eureka fake.
	eureka.Close()

	// Segunda resolução: Eureka fora do ar, mas ainda dentro do TTL,
	// então nem deveria tentar consultar — deve retornar o valor já
	// cacheado sem erro.
	secondURL, err := discovery.ResolveURL(context.Background(), "currency-service")
	if err != nil {
		t.Fatalf("expected cached url to be served without error, got: %v", err)
	}
	if secondURL != firstURL {
		t.Fatalf("expected cached url %q, got %q", firstURL, secondURL)
	}
}

func TestHTTPCurrencyClient_NotFoundDoesNotTripCircuitBreaker(t *testing.T) {
	callCount := 0
	currencyServer := newFakeCurrencyService(t, func(w http.ResponseWriter, r *http.Request) {
		callCount++
		// As primeiras 6 chamadas (acima do threshold de 5 do ReadyToTrip)
		// retornam 404. Se isso contasse como falha de circuito, a 7ª
		// chamada já encontraria o circuito aberto.
		if callCount <= 6 {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(CurrencyResponse{ConversionRate: 0.2})
	})
	eureka := newFakeEureka(t, currencyServer)

	discovery := NewServiceDiscovery(eureka.URL, &http.Client{Timeout: 2 * time.Second}, 10*time.Second)
	client := NewHTTPCurrencyClient(discovery, 2*time.Second)

	for i := 0; i < 6; i++ {
		_, err := client.GetCurrency(context.Background(), "BRL", "XYZ")
		if !errors.Is(err, ErrCurrencyNotFound) {
			t.Fatalf("call %d: expected ErrCurrencyNotFound, got: %v", i, err)
		}
	}

	// Depois de 6 erros de negócio, o circuito deveria continuar fechado:
	// esta chamada bem-sucedida não deve ser bloqueada por ErrOpenState.
	result, err := client.GetCurrency(context.Background(), "BRL", "USD")
	if err != nil {
		t.Fatalf("expected circuit to remain closed after business errors, got: %v", err)
	}
	if result.ConversionRate != 0.2 {
		t.Fatalf("expected rate 0.2, got %v", result.ConversionRate)
	}
}
