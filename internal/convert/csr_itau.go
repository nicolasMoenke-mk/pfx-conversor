package convert

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
)

type ItauCSR struct {
	CommonName         string
	Organization       string
	OrganizationalUnit string
	Country            string
	State              string
	Locality           string
	Email              string
}

func GenNewCSR(dados ItauCSR) (csrPEM []byte, keyPEM []byte, err error) {

	if dados.CommonName == "" {
		return nil, nil, errors.New("Common Name is required to generated the CSR")
	}

	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, errors.New("Failed to generated RSA Key")
	}

	subject := pkix.Name{
		CommonName:         dados.CommonName,
		Organization:       noNull(dados.Organization),
		OrganizationalUnit: noNull(dados.OrganizationalUnit),
		Country:            noNull(dados.Country),
		Province:           noNull(dados.State),
		Locality:           noNull(dados.Locality),
	}

	template := x509.CertificateRequest{
		Subject:            subject,
		SignatureAlgorithm: x509.SHA256WithRSA,
	}

	if dados.Email != "" {
		template.EmailAddresses = []string{dados.Email}
	}

	csrDer, err := x509.CreateCertificateRequest(rand.Reader, &template, privKey)
	if err != nil {
		return nil, nil, errors.New("Failed to generated CSR")
	}

	csrPEM = pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csrDer,
	})

	keyDer, err := x509.MarshalPKCS8PrivateKey(privKey)
	if err != nil {
		return nil, nil, errors.New("Failed to convert the key to PKCS8 format")
	}

	keyPEM = pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: keyDer,
	})

	return csrPEM, keyPEM, nil

}

/*
Função que determina se um elemento vier vazio, faz o tratamento ignorar e omitir
o campo no Subject da struct, retornando um slice de string devido a biblioteca pkix,
os elementos da struct Name da lib retornam este tipo
*/
func noNull(value string) []string {
	if value == "" {
		return nil
	}

	return []string{value}
}
