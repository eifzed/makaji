package mongodb

type Config struct {
	URI                  string `json:"uri"`
	EnableSSL            bool   `json:"enable_ssl"`
	CertificateWritePath string `json:"certificate_write_path"`
	ClientKey            string `json:"client_key"`
	ClientSecret         string `json:"client_secret"`
	ServiceCA            string `json:"server_ca"`
}
