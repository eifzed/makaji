package jwt

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"
	"time"
)

type JWTCertificate struct {
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
}

type RouteRoles struct {
	Roles []Role `yaml:"roles"`
}

type Role struct {
	ID   int64  `yaml:"id" json:"id" xorm:"role_id"`
	Name string `yaml:"name" json:"name" xorm:"role_name"`
}

type JWTPayload struct {
	UserID         int64
	Name           string
	Email          string
	Username       string
	Roles          []Role
	PasswordHashed string
	GeneratedUnix  int64
	ExpiredUnix    int64
}

type JWTHeader struct {
	Algorithm string
	Type      string
}

var (
	ErrMarshal   = errors.New("JSON marshal failed")
	ErrSigning   = errors.New("token signing failed")
	ErrInvalid   = errors.New("invalid token")
	ErrExpired   = errors.New("token is expired")
	ErrForbidden = errors.New("user is forbidden to access this resource")
)

func GenerateToken(payload JWTPayload, privateKey string, expiredAfterMinutes int64) (string, error) {
	now := time.Now()
	payload.GeneratedUnix = now.Unix()
	payload.ExpiredUnix = now.Add(time.Duration(expiredAfterMinutes) * time.Minute).Unix()

	header := JWTHeader{
		Algorithm: "RS256",
		Type:      "JWT",
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", ErrMarshal
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "", ErrMarshal
	}
	header64 := base64.StdEncoding.EncodeToString(headerJSON)
	payload64 := base64.StdEncoding.EncodeToString(payloadJSON)
	data := fmt.Sprintf("%s.%s", header64, payload64)
	parsedPrivateKey, err := parsePrivateKey([]byte(privateKey))
	if err != nil {
		return "", err
	}
	h := sha256.New()
	h.Write([]byte(data))
	d := h.Sum(nil)

	signedData, err := rsa.SignPKCS1v15(rand.Reader, parsedPrivateKey, crypto.SHA256, d)
	if err != nil {
		return "", ErrSigning
	}
	signature64 := base64.StdEncoding.EncodeToString(signedData)
	return fmt.Sprintf("%s.%s.%s", header64, payload64, signature64), nil
}

func parsePrivateKey(pemBytes []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("ssh: no key found")
	}
	switch block.Type {
	case "RSA PRIVATE KEY":
		return x509.ParsePKCS1PrivateKey(block.Bytes)
	default:
		return nil, fmt.Errorf("ssh: unsupported key type %q", block.Type)
	}
}

func parsePublicKey(pemBytes []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("ssh: no key found")
	}
	switch block.Type {
	case "PUBLIC KEY":
		parser, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, errors.New("failed to get public parser")
		}
		return parser.(*rsa.PublicKey), nil
	default:
		return nil, fmt.Errorf("ssh: unsupported key type %q", block.Type)
	}
}

func DecodeToken(token string, publicKey string) (*JWTPayload, error) {
	payloadJSON, err := getJWTPayload(token, publicKey)
	if err != nil {
		return nil, err
	}
	payload := &JWTPayload{}
	err = json.Unmarshal(payloadJSON, payload)
	if err != nil {
		return nil, ErrInvalid
	}
	if payload.ExpiredUnix < time.Now().Unix() {
		return nil, ErrExpired
	}
	return payload, nil
}

func getJWTPayload(token string, publicKey string) ([]byte, error) {
	tokenList := strings.Split(token, ".")
	if len(tokenList) != 3 {
		return nil, ErrInvalid
	}
	err := verifyJWTToken(fmt.Sprintf("%s.%s", tokenList[0], tokenList[1]), tokenList[2], publicKey)
	if err != nil {
		return nil, ErrInvalid
	}
	payloadJSON, err := base64.StdEncoding.DecodeString(tokenList[1])
	if err != nil {
		return nil, ErrInvalid
	}
	return payloadJSON, nil

}

func verifyJWTToken(data string, signature string, publicKey string) error {
	signed, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return fmt.Errorf("signature is not base64: %v", err)
	}
	publicParsed, err := parsePublicKey([]byte(publicKey))
	if err != nil {
		return ErrInvalid
	}
	h := sha256.New()
	h.Write([]byte(data))
	d := h.Sum(nil)

	return rsa.VerifyPKCS1v15(publicParsed, crypto.SHA256, d, signed)
}
