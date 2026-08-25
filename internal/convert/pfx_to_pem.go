package convert

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"errors"

	"golang.org/x/crypto/pkcs12"
)

func decodePFX(pfxData []byte, password string) (crypto.PrivateKey, *x509.Certificate, error) {

	key, cert, err := pkcs12.Decode(pfxData, password)
	if err != nil {
		return nil, nil, errors.New("Falha ao ler certificado!")
	}

	if key == nil {
		return nil, nil, errors.New("Falha ao carregar chave privada!")
	}

	return key, cert, nil
}

func privateKeyToPEM(key crypto.PrivateKey) ([]byte, error) {

	derBytes, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return nil, errors.New("Certificado sem Chave Privada válida")
	}

	block := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: derBytes,
	}

	return pem.EncodeToMemory(block), nil
}

func certificateToPEM(cert *x509.Certificate) ([]byte, error) {

	block := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	}

	return pem.EncodeToMemory(block), nil
}

func ConvertPFX(pfxData []byte, password string) ([]byte, []byte, error) {

	privKey, cert, err := decodePFX(pfxData, password)
	if err != nil {
		return nil, nil, errors.New("Falha ao abrir o pfx, verifique se a senha está correta!")
	}

	keyPEM, err := privateKeyToPEM(privKey)
	if err != nil {
		return nil, nil, errors.New("Falha ao criar a Chave Privada")
	}

	certPEM, err := certificateToPEM(cert)
	if err != nil {
		return nil, nil, errors.New("Falha ao criar a Chave Pública")
	}

	return keyPEM, certPEM, nil
}
