package clients

// EurekaAppResponse é a resposta de GET {eurekaURL}/apps/{APP_NAME}.
type EurekaAppResponse struct {
	Application EurekaApplication `json:"application"`
}

type EurekaApplication struct {
	Name      string                `json:"name"`
	Instances []EurekaInstanceQuery `json:"instance"`
}

// EurekaInstanceQuery é o formato de instância retornado pela consulta —
// diferente do InstanceInfo usado no registro. Repare que aqui o "port"
// vem como número (json:"$" int), não string como no payload que enviamos
// ao registrar; é uma assimetria do próprio Eureka entre request e response.
type EurekaInstanceQuery struct {
	InstanceID string          `json:"instanceId"`
	HostName   string          `json:"hostName"`
	App        string          `json:"app"`
	IPAddr     string          `json:"ipAddr"`
	Status     string          `json:"status"`
	Port       EurekaPortQuery `json:"port"`
}

type EurekaPortQuery struct {
	Number  int    `json:"$"`
	Enabled string `json:"@enabled"`
}
