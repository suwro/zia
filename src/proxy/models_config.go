package proxy

type ConfigStruct struct {
	Service Service    `json:"service"`
	Domains []ziaProxy `json:"domains"`
}

type Service struct {
	IP       string `json:"ip" validate:"omitempty,ipv4"`
	Port     int    `json:"port" validate:"required,numeric"`
	HostName string `json:"hostname" validate:"required,hostname"`
	SSL      bool   `json:"ssl" validate:"boolean"`
}

type ziaProxy struct {
	Host   string `json:"host" validate:"required"`
	Target string `json:"target" validate:"required,url"`
}
