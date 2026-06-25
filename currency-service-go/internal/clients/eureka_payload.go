package clients

// InstanceRegistration é o payload enviado no POST /eureka/apps/{APP_NAME}.
// Espelha a estrutura com.netflix.appinfo.InstanceInfo serializada como
// JSON pelo Eureka Server. O campo raiz "instance" é obrigatório no formato
// que o Eureka aceita.
type InstanceRegistration struct {
	Instance InstanceInfo `json:"instance"`
}

type InstanceInfo struct {
	InstanceID       string         `json:"instanceId"`
	HostName         string         `json:"hostName"`
	App              string         `json:"app"`
	IPAddr           string         `json:"ipAddr"`
	Status           string         `json:"status"`
	Port             PortInfo       `json:"port"`
	SecurePort       PortInfo       `json:"securePort"`
	HomePageURL      string         `json:"homePageUrl"`
	StatusPageURL    string         `json:"statusPageUrl"`
	HealthCheckURL   string         `json:"healthCheckUrl"`
	VipAddress       string         `json:"vipAddress"`
	SecureVipAddress string         `json:"secureVipAddress"`
	DataCenterInfo   DataCenterInfo `json:"dataCenterInfo"`
}

// PortInfo precisa desse formato específico porque o Eureka serializa
// number como string e enabled como string "true"/"false" no XML
// original — o JSON do Eureka mantém essa peculiaridade por compatibilidade.
type PortInfo struct {
	Number  string `json:"$"`
	Enabled string `json:"@enabled"`
}

type DataCenterInfo struct {
	Name      string `json:"name"`
	ClassName string `json:"@class"`
}
