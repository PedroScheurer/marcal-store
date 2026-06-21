package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

type cachedResolution struct {
	url       string
	expiresAt time.Time
}

// ServiceDiscovery resolve, via Eureka, o endereço atual de uma aplicação
// registrada, com cache local por TTL — equivalente ao comportamento do
// Spring Cloud LoadBalancer, que atualiza periodicamente sua lista de
// instâncias em vez de consultar o Eureka a cada chamada.
type ServiceDiscovery struct {
	eurekaURL  string
	httpClient *http.Client
	ttl        time.Duration

	mu    sync.Mutex
	cache map[string]cachedResolution
}

func NewServiceDiscovery(eurekaURL string, httpClient *http.Client, ttl time.Duration) *ServiceDiscovery {
	return &ServiceDiscovery{
		eurekaURL:  eurekaURL,
		httpClient: httpClient,
		ttl:        ttl,
		cache:      make(map[string]cachedResolution),
	}
}

// ResolveURL retorna a base URL HTTP de uma instância UP de appName.
//
// Se houver uma entrada em cache ainda válida (dentro do TTL), retorna
// direto sem consultar o Eureka. Caso contrário, consulta o Eureka; se
// a consulta falhar mas existir uma entrada antiga em cache (mesmo
// expirada), usa essa entrada como fallback — preferimos servir um
// endereço "potencialmente desatualizado" a falhar completamente quando
// o Eureka está temporariamente fora do ar.
func (d *ServiceDiscovery) ResolveURL(ctx context.Context, appName string) (string, error) {
	key := strings.ToUpper(appName)

	d.mu.Lock()
	cached, found := d.cache[key]
	d.mu.Unlock()

	if found && time.Now().Before(cached.expiresAt) {
		return cached.url, nil
	}

	url, err := d.fetchFromEureka(ctx, key)
	if err != nil {
		if found {
			// Eureka indisponível, mas temos um valor antigo: melhor
			// tentar com ele do que falhar a chamada inteira.
			return cached.url, nil
		}
		return "", err
	}

	d.mu.Lock()
	d.cache[key] = cachedResolution{url: url, expiresAt: time.Now().Add(d.ttl)}
	d.mu.Unlock()

	return url, nil
}

func (d *ServiceDiscovery) fetchFromEureka(ctx context.Context, appName string) (string, error) {
	url := fmt.Sprintf("%s/apps/%s", d.eurekaURL, appName)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("create eureka discovery request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("call eureka discovery: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("eureka discovery failed for %s: status %d", appName, resp.StatusCode)
	}

	var appResponse EurekaAppResponse
	if err := json.NewDecoder(resp.Body).Decode(&appResponse); err != nil {
		return "", fmt.Errorf("decode eureka discovery response: %w", err)
	}

	for _, instance := range appResponse.Application.Instances {
		if instance.Status == "UP" {
			return fmt.Sprintf("http://%s:%d", instance.HostName, instance.Port.Number), nil
		}
	}

	return "", fmt.Errorf("no UP instance found for app %s", appName)
}
