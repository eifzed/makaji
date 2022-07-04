package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/eifzed/joona/lib/common"
	"github.com/eifzed/joona/lib/database/mongodb"
	"github.com/eifzed/joona/lib/utility/jwt"
)

type SecretVault struct {
	Data     *DataVault `json:"data"`
	Metadata *Metadata  `json:"metadata"`
}

type DataVault struct {
	MongoDBConfig  *mongodb.Config     `json:"mongo_db_config"`
	JWTCertificate *jwt.JWTCertificate `json:"jwt_certificate"`
}

type Metadata struct {
	CreatedTime  string  `json:"created_time"`
	Destroyed    bool    `json:"destroyed"`
	Version      float32 `json:"version"`
	DeletionTime string  `json:"deletion_time"`
}

func GetSecrets() *SecretVault {
	env := "production"
	vaultPath := "/etc/joona-secret/"

	if common.IsDevelopment() {
		dir, _ := os.Getwd()
		env = "development"
		vaultPath = dir + "/files" + "/etc/joona-secret/"
	}
	vaultPath = vaultPath + "joona-secret" + "." + env + ".json"
	vaultFile, err := os.Open(vaultPath)
	if err != nil {
		log.Fatalln("Path fault not found:", err)
	}
	configByte, err := ioutil.ReadAll(vaultFile)
	if err != nil {
		log.Fatalln("Path fault not found:", err)
	}
	cfgVault := SecretVault{}
	err = json.Unmarshal(configByte, &cfgVault)
	if err != nil {
		log.Fatalln("Failed get vault config:", err)
	}
	if cfgVault.Data == nil {
		log.Fatalln("Failed config vault nil on data")
	}
	return &cfgVault
}
