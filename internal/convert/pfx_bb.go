package convert

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"errors"

	"golang.org/x/crypto/pkcs12"
	sslpkcs12 "software.sslmate.com/src/go-pkcs12"
)

type CadeiaCertificado struct {
	Raiz          string `json:"raiz,omitempty"`
	NomeRaiz      string `json:"nomeRaiz,omitempty"`
	Intermed1     string `json:"intermed1,omitempty"`
	NomeIntermed1 string `json:"nomeIntermed1,omitempty"`
	Intermed2     string `json:"intermed2,omitempty"`
	NomeIntermed2 string `json:"nomeIntermed2,omitempty"`
}

func decodeBBPFX(pfxData []byte, password string) (crypto.PrivateKey, *x509.Certificate, error) {
	key, cert, err := pkcs12.Decode(pfxData, password)
	if err != nil {
		return nil, nil, errors.New("Falha ao ler conteúdo do certificado")
	}

	if key == nil {
		return nil, nil, errors.New("Falha ao carregar Chave Privada")
	}

	return key, cert, nil
}

func privateKeyBBtoPEM(key crypto.PrivateKey) ([]byte, error) {
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

func certificateBBtoPEM(cert *x509.Certificate) ([]byte, error) {

	block := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	}

	return pem.EncodeToMemory(block), nil
}

func genKeysBB(pfxData []byte, password string) ([]byte, []byte, error) {

	privKey, cert, err := decodeBBPFX(pfxData, password)
	if err != nil {
		return nil, nil, errors.New("Falha ao carregar o pfx")
	}

	keyPEM, err := privateKeyBBtoPEM(privKey)
	if err != nil {
		return nil, nil, errors.New("Erro ao tentar extrair a Chave Privada")
	}

	certPEM, err := certificateBBtoPEM(cert)
	if err != nil {
		return nil, nil, errors.New("Erro ao tentar extrair a Chave Pública")
	}

	return keyPEM, certPEM, nil

}

func extrairCadeiaBB(pfxData []byte, password string) (*CadeiaCertificado, error) {
	_, _, chain, err := sslpkcs12.DecodeChain(pfxData, password)
	if err != nil {
		return nil, errors.New("Falha ao abrir o pfx, tente novamente")
	}

	res := CadeiaCertificado{}

	for _, cert := range chain {
		pemBytes, err := certificateToPEM(cert)
		if err != nil {
			continue
		}
		nome := cert.Subject.CommonName
		conteudo := string(pemBytes)

		if cert.Issuer.CommonName == cert.Subject.CommonName {
			res.Raiz = conteudo
			res.NomeRaiz = nome
		} else if res.Intermed2 == "" {
			res.Intermed2 = conteudo
			res.NomeIntermed2 = nome
		} else {
			res.Intermed1 = conteudo
			res.NomeIntermed1 = nome
		}
	}

	return &res, nil
}
