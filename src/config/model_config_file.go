package config

type Config struct {
	Port         int         `json:"port"`
	Ip           string      `json:"ip"`
	DomainName   string      `json:"domain_name"`
	SSL          bool        `json:"ssl"`
	Cert         string      `json:"cert"`
	Key          string      `json:"key"`
	Timeout      int         `json:"timeout"`
	NoAccessPage string      `json:"no_access_page"`
	Targets      []string    `json:"targets"`
	AllowedIPs   []AllowedIP `json:"allowed_ips"`
}

type AllowedIP struct {
	Ip   string `json:"ip"`
	Name string `json:"name"`
}
