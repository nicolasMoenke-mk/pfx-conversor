package convert

import (
	"crypto/x509"
	"errors"

	sslpkcs12 "software.sslmate.com/src/go-pkcs12"
)

// ConvertPfxNext usa sslpkcs12.DecodeChain pelo mesmo motivo de decodePFX em
// pfx_to_pem.go: um .pfx com a cadeia embutida quebrava a decodificação
// antes mesmo de chegar a reencodar para o Next. A cadeia (se houver) é
// descartada de propósito — o Next só usa o par chave/certificado da folha,
// já era assim mesmo antes desta correção (Encode recebe []*x509.Certificate{}).
func ConvertPfxNext(pfxData []byte, senhaAtual, novaSenha string) ([]byte, error) {
	key, cert, _, err := sslpkcs12.DecodeChain(pfxData, senhaAtual)
	if err != nil {
		return nil, errors.New("Falha ao abrir o certificado, verifique se a senha está correta")
	}

	newPfx, err := sslpkcs12.Modern.Encode(key, cert, []*x509.Certificate{}, novaSenha)
	if err != nil {
		return nil, errors.New("Erro inesperado ao tentar converter o certificado")
	}

	return newPfx, nil
}
