package config

type Configuration struct {
	Neo4j           Neo4jConfiguration `json:"neo4j"`
	Mail            MailConfiguration  `json:"mail"`
	FrontendBaseUrl string             `json:"frontend_base_url"`
}

type Neo4jConfiguration struct {
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type MailConfiguration struct {
	Server      string `json:"server"`
	ServerPort  int    `json:"server_port"`
	Password    string `json:"password"`
	Username    string `json:"username"`
	MailAddress string `json:"mail_address"`
}


