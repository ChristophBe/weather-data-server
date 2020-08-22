package config

type Configuration struct {
	Neo4j           Neo4jConfiguration `json:"neo4j"`
	Mail            MailConfiguration  `json:"mail"`
	Auth            AuthConfiguration  `json:"auth"`
	RSAKeyFile      string             `json:"rsa_key_file"`
	FrontendBaseUrl string             `json:"frontend_base_url"`
}

type Neo4jConfiguration struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}
type AuthConfiguration struct {
	AuthKey string `json:"auth_key"`
}

type MailConfiguration struct {
	Server      string `json:"server"`
	ServerPort  int    `json:"server_port"`
	Password    string `json:"password"`
	Username    string `json:"username"`
	MailAddress string `json:"mail_address"`
	MailName string `json:"mail_name"`
}
