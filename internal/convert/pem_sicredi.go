package convert

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
)

// decodeCertificateBytes aceita tanto PEM quanto DER "cru" — o .cer que o
// Sicredi entrega já costuma ser um texto PEM (apesar da extensão .cer),
// mas alguns navegadores/portais baixam a variante binária DER.
func decodeCertificateBytes(data []byte) (*x509.Certificate, error) {
	if block, _ := pem.Decode(data); block != nil {
		return x509.ParseCertificate(block.Bytes)
	}
	return x509.ParseCertificate(data)
}

// ConvertCERtoPEM normaliza o certificado (.cer) que o Sicredi entrega para
// PEM — o mesmo resultado de "openssl x509 -in cert.cer -out cert.pem".
func ConvertCERtoPEM(cerData []byte) ([]byte, error) {
	cert, err := decodeCertificateBytes(cerData)
	if err != nil {
		return nil, errors.New("Falha ao ler o certificado .cer — verifique se o arquivo não está corrompido")
	}
	return certificateToPEM(cert)
}

// ConvertKEYtoPEM normaliza a chave privada RSA (.key) que o Sicredi entrega
// para PEM — equivalente a "openssl rsa -in chave.key -out chave.pem".
// Aceita tanto PKCS1 ("RSA PRIVATE KEY") quanto PKCS8 ("PRIVATE KEY"), em
// PEM ou DER, e sempre devolve PKCS8 — o mesmo formato usado nas demais
// conversões desta ferramenta.
func ConvertKEYtoPEM(keyData []byte) ([]byte, error) {
	der := keyData
	if block, _ := pem.Decode(keyData); block != nil {
		der = block.Bytes
	}

	if key, err := x509.ParsePKCS1PrivateKey(der); err == nil {
		return privateKeyToPEM(key)
	}

	if key, err := x509.ParsePKCS8PrivateKey(der); err == nil {
		return privateKeyToPEM(key)
	}

	return nil, errors.New("Falha ao ler a chave privada .key — formato não reconhecido (esperado RSA)")
}
