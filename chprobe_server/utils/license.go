package utils

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
)

var publicKey *rsa.PublicKey

func InitRSA() error {
	publicKeyBytes, err := os.ReadFile("./utils/public.pem")
	if err != nil {
		return fmt.Errorf("failed to read public key file: %w", err)
	}

	block, _ := pem.Decode(publicKeyBytes)
	if block == nil {
		return fmt.Errorf("failed to parse PEM block")
	}

	publicKeyParsed, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse public key: %w", err)
	}

	publicKeyRSA, ok := publicKeyParsed.(*rsa.PublicKey)
	if !ok {
		return fmt.Errorf("not an RSA public key")
	}

	publicKey = publicKeyRSA

	return nil
}

func VerifyLicenseString(licenseString string) (map[string]interface{}, bool, error) {
	InitRSA()
	if publicKey == nil {
		return nil, false, fmt.Errorf("RSA public key not initialized")
	}

	payloadBytes, err := base64.StdEncoding.DecodeString(licenseString)
	if err != nil {
		return nil, false, err
	}

	var licensePayload map[string]interface{}
	if err := json.Unmarshal(payloadBytes, &licensePayload); err != nil {
		return nil, false, err
	}

	data, ok := licensePayload["data"].(string)
	if !ok {
		return nil, false, fmt.Errorf("invalid license payload: missing data")
	}

	signatureStr, ok := licensePayload["signature"].(string)
	if !ok {
		return nil, false, fmt.Errorf("invalid license payload: missing signature")
	}

	signature, err := base64.StdEncoding.DecodeString(signatureStr)
	if err != nil {
		return nil, false, err
	}

	jsonData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, false, err
	}

	hashed := sha256.Sum256(jsonData)

	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], signature)
	if err != nil {
		return nil, false, err
	}

	var licenseData map[string]interface{}
	if err := json.Unmarshal(jsonData, &licenseData); err != nil {
		return nil, false, err
	}

	return licenseData, true, nil
}
