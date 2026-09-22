package convert

import (
	"bytes"
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"errors"

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

// autoAssinado reconhece um certificado raiz pela definição real de
// self-signed (emissor e titular são o mesmo DN por completo) — comparar só
// o Common Name, como o código fazia antes, falha sempre que a autoridade
// certificadora usa um CN vazio ou formatado de um jeito que não bate com o
// do titular, mesmo sendo o mesmo certificado.
func autoAssinado(cert *x509.Certificate) bool {
	return bytes.Equal(cert.RawIssuer, cert.RawSubject)
}

// ordenarCadeiaBB monta o caminho de certificação da folha até a raiz
// verificando de fato quem assinou quem (CheckSignatureFrom), em vez de
// assumir a ordem em que os certificados apareceram dentro do .pfx — essa
// ordem não é garantida, e um .pfx que traga algum certificado extra (ex.:
// uma AC intermediária antiga, ainda válida, de uma transição de cadeia)
// podia antes tomar o lugar do certificado certo e deixar a cadeia
// incompleta. Certificados que não fazem parte do caminho real da folha até
// a raiz são ignorados.
func ordenarCadeiaBB(folha *x509.Certificate, pool []*x509.Certificate) []*x509.Certificate {

	var caminho []*x509.Certificate

	atual := folha
	restantes := append([]*x509.Certificate{}, pool...)

	for !autoAssinado(atual) {

		idx := -1
		for i, candidato := range restantes {
			if atual.CheckSignatureFrom(candidato) == nil {
				idx = i
				break
			}
		}

		if idx == -1 {
			break
		}

		caminho = append(caminho, restantes[idx])
		atual = restantes[idx]
		restantes = append(restantes[:idx], restantes[idx+1:]...)

	}

	return caminho
}

// ConvertPFXBB extrai, a partir do .pfx do cliente, tudo que o Banco do
// Brasil exige para o mTLS: a chave pública/privada da aplicação e a cadeia
// de certificação (raiz + intermediários). É o único banco cujo portal
// Developer pede a cadeia completa, além do par de chaves comum aos demais.
//
// Usa sslpkcs12.DecodeChain (não golang.org/x/crypto/pkcs12, usado no
// conversor genérico) porque o .pfx do BB normalmente já vem com a cadeia
// inteira dentro do arquivo — e o pacote x/crypto só sabe ler um PFX com
// exatamente um par chave/certificado, falhando ("expected exactly two safe
// bags") assim que há mais de um certificado no pacote.
func ConvertPFXBB(pfxData []byte, password string) ([]byte, []byte, *CadeiaCertificado, error) {

	privKey, cert, chain, err := sslpkcs12.DecodeChain(pfxData, password)
	if err != nil {
		return nil, nil, nil, errors.New("Falha ao abrir o pfx, verifique se a senha está correta!")
	}

	keyPEM, err := privateKeyBBtoPEM(privKey)
	if err != nil {
		return nil, nil, nil, errors.New("Erro ao tentar extrair a Chave Privada")
	}

	certPEM, err := certificateBBtoPEM(cert)
	if err != nil {
		return nil, nil, nil, errors.New("Erro ao tentar extrair a Chave Pública")
	}

	cadeia := CadeiaCertificado{}

	caminho := ordenarCadeiaBB(cert, chain)

	for i, c := range caminho {

		pemBytes, err := certificateBBtoPEM(c)
		if err != nil {
			continue
		}

		nome := c.Subject.CommonName
		conteudo := string(pemBytes)

		switch {
		case autoAssinado(c):
			cadeia.Raiz = conteudo
			cadeia.NomeRaiz = nome
		case i == 0:
			// primeiro da folha para cima: o intermediário que assina
			// diretamente o certificado da aplicação.
			cadeia.Intermed2 = conteudo
			cadeia.NomeIntermed2 = nome
		default:
			cadeia.Intermed1 = conteudo
			cadeia.NomeIntermed1 = nome
		}

	}

	return keyPEM, certPEM, &cadeia, nil
}
