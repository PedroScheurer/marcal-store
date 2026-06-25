package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// EurekaClient é o equivalente Go ao comportamento automático do
// Spring Cloud Netflix Eureka Client (@EnableDiscoveryClient).
// Cobre as três operações do protocolo REST do Eureka: registro,
// heartbeat (renovação de lease) e deregistro.
type EurekaClient struct {
	eurekaURL  string // ex.: "http://localhost:8761/eureka"
	appName    string
	instanceID string
	hostName   string
	port       string
	httpClient *http.Client
}

// NewEurekaClient cria um client do Eureka para esta instância do serviço.
// instanceID segue a convenção "{hostName}:{appName}:{port}", que é o
// mesmo padrão usado pelo Spring Cloud Netflix Eureka Client por padrão.
func NewEurekaClient(eurekaURL, appName, hostName, port string) *EurekaClient {
	instanceID := fmt.Sprintf("%s:%s:%s", hostName, appName, port)

	return &EurekaClient{
		eurekaURL:  eurekaURL,
		appName:    appName,
		instanceID: instanceID,
		hostName:   hostName,
		port:       port,
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
}

// Register é o equivalente ao POST /eureka/apps/{APP_NAME} feito
// automaticamente pelo Spring no startup da aplicação.
func (c *EurekaClient) Register(ctx context.Context) error {
	baseURL := fmt.Sprintf("http://%s:%s", c.hostName, c.port)

	registration := &InstanceRegistration{
		Instance: InstanceInfo{
			InstanceID: c.instanceID,
			HostName:   c.hostName,
			App:        c.appName,
			IPAddr:     c.hostName,
			Status:     "UP",
			Port: PortInfo{
				Number:  c.port,
				Enabled: "true",
			},
			SecurePort: PortInfo{
				Number:  "443",
				Enabled: "false",
			},
			HomePageURL:      baseURL + "/",
			StatusPageURL:    baseURL + "/management/info",
			HealthCheckURL:   baseURL + "/management/health",
			VipAddress:       c.appName,
			SecureVipAddress: c.appName,
			DataCenterInfo: DataCenterInfo{
				Name:      "MyOwn",
				ClassName: "com.netflix.appinfo.MyDataCenterInfo",
			},
		},
	}

	body, err := json.Marshal(registration)
	if err != nil {
		return fmt.Errorf("marshal eureka registration: %w", err)
	}

	url := fmt.Sprintf("%s/apps/%s", c.eurekaURL, c.appName)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create eureka register request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("call eureka register: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("eureka register failed: status %d", resp.StatusCode)
	}

	return nil
}

// Heartbeat é o equivalente ao PUT /eureka/apps/{APP_NAME}/{INSTANCE_ID}
// que o Spring envia periodicamente (lease-renewal-interval-in-seconds,
// 30s por padrão) para renovar o "lease" da instância no Eureka.
//
// Se o Eureka responder 404, significa que a instância expirou por lá
// (não recebeu heartbeat dentro do timeout de eviction) e precisa ser
// registrada novamente — é por isso que o chamador (StartHeartbeatLoop)
// trata esse caso re-chamando Register.
func (c *EurekaClient) Heartbeat(ctx context.Context) error {
	url := fmt.Sprintf("%s/apps/%s/%s", c.eurekaURL, c.appName, c.instanceID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, nil)
	if err != nil {
		return fmt.Errorf("create eureka heartbeat request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("call eureka heartbeat: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return errEurekaInstanceNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("eureka heartbeat failed: status %d", resp.StatusCode)
	}

	return nil
}

// Deregister é o equivalente ao DELETE /eureka/apps/{APP_NAME}/{INSTANCE_ID}
// que o Spring envia no shutdown gracioso da aplicação (via shutdown hook),
// removendo a instância do Eureka imediatamente em vez de esperar o
// timeout de eviction.
func (c *EurekaClient) Deregister(ctx context.Context) error {
	url := fmt.Sprintf("%s/apps/%s/%s", c.eurekaURL, c.appName, c.instanceID)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("create eureka deregister request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("call eureka deregister: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("eureka deregister failed: status %d", resp.StatusCode)
	}

	return nil
}

// StartHeartbeatLoop envia heartbeats periodicamente até que ctx seja
// cancelado, equivalente ao agendador interno do Spring Cloud Netflix
// Eureka Client. Deve ser chamada em uma goroutine separada (go
// client.StartHeartbeatLoop(ctx, 30*time.Second)), já que ela bloqueia
// até o contexto ser cancelado.
func (c *EurekaClient) StartHeartbeatLoop(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := c.Heartbeat(ctx)
			if err == nil {
				continue
			}

			if err == errEurekaInstanceNotFound {
				log.Println("eureka heartbeat: instance not found, re-registering")
				if regErr := c.Register(ctx); regErr != nil {
					log.Printf("eureka re-register failed: %v", regErr)
				}
				continue
			}

			log.Printf("eureka heartbeat failed: %v", err)
		}
	}
}

// errEurekaInstanceNotFound sinaliza que o Eureka respondeu 404 ao
// heartbeat, indicando que a instância expirou por lá e precisa se
// re-registrar antes de continuar enviando heartbeats.
var errEurekaInstanceNotFound = fmt.Errorf("eureka: instance not registered")
