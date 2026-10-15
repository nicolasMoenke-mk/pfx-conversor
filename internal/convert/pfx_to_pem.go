package convert

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"errors"

	sslpkcs12 "software.sslmate.com/src/go-pkcs12"
)

// decodePFX usa sslpkcs12.DecodeChain (não golang.org/x/crypto/pkcs12) pelo
// mesmo motivo do BB: certificados brasileiros costumam vir com a cadeia
// (raiz/intermediários) já embutida no .pfx, e o pacote x/crypto só sabe
// ler um PFX com exatamente um par chave/certificado — falha com "expected
// exactly two safe bags" assim que há mais de um certificado no arquivo.
// Aqui só a folha (chave + certificado da aplicação) importa, então a
// cadeia retornada é descartada.
func decodePFX(pfxData []byte, password string) (crypto.PrivateKey, *x509.Certificate, error) {

	key, cert, _, err := sslpkcs12.DecodeChain(pfxData, password)
	if err != nil {
		return nil, nil, errors.New("Failed to load certificate")
	}

	if key == nil {
		return nil, nil, errors.New("Failed to load private key")
	}

	return key, cert, nil
}

func privateKeyToPEM(key crypto.PrivateKey) ([]byte, error) {

	derBytes, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return nil, errors.New("Certificate without a valid private key")
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
		return nil, nil, errors.New("Failed to load pfx, check if password is correct")
	}

	keyPEM, err := privateKeyToPEM(privKey)
	if err != nil {
		return nil, nil, errors.New("Failed to generated the private key")
	}

	certPEM, err := certificateToPEM(cert)
	if err != nil {
		return nil, nil, errors.New("Failed to generated the public key")
	}

	return keyPEM, certPEM, nil
}
