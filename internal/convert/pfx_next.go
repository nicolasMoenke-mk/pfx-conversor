package convert

import (
	"crypto/x509"
	"errors"

	"golang.org/x/crypto/pkcs12"
	sslpkcs12 "software.sslmate.com/src/go-pkcs12"
)

func ConvertPfxNext(pfxData []byte, senhaAtual, novaSenha string) ([]byte, error) {
	key, cert, err := pkcs12.Decode(pfxData, senhaAtual)
	if err != nil {
		return nil, errors.New("Senha inserida incorreta, tente novamente")
	}

	newPfx, err := sslpkcs12.Modern.Encode(key, cert, []*x509.Certificate{}, novaSenha)
	if err != nil {
		return nil, errors.New("Erro inesperado ao tentar converter o certificado")
	}

	return newPfx, nil
}
