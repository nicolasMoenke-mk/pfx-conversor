// Package remessa valida a integridade estrutural de arquivos de remessa
// bancária nos layouts CNAB 240 e CNAB 400: tamanho de linha, sequência dos
// registros e a ordem dos tipos de registro (header/detalhe/trailer).
//
// O objetivo não é validar o conteúdo de negócio de cada campo (valores,
// datas, contas etc.), e sim encontrar rapidamente onde o arquivo está
// estruturalmente quebrado — normalmente sinal de uma linha faltando,
// duplicada ou fora de ordem — o que é praticamente impossível de achar a
// olho nu em arquivos com milhares de linhas.
package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type TipoCnab string

const (
	CNAB240 TipoCnab = "CNAB240"
	CNAB400 TipoCnab = "CNAB400"
)

type Erro struct {
	Linha     int    `json:"linha"`
	Categoria string `json:"categoria"`
	Descricao string `json:"descricao"`
}

type Resultado struct {
	TipoCnab        TipoCnab `json:"tipoCnab"`
	TotalLinhas     int      `json:"totalLinhas"`
	QuantidadeLotes int      `json:"quantidadeLotes,omitempty"`
	Erros           []Erro   `json:"erros"`
	Valido          bool     `json:"valido"`
}

// splitLinhas separa o conteúdo em linhas, aceitando tanto quebra de linha
// Unix (\n) quanto Windows (\r\n), e descarta uma eventual linha vazia final
// causada por o arquivo terminar com quebra de linha. Espaços internos das
// linhas NUNCA são removidos: em CNAB o preenchimento com espaço faz parte
// do layout de largura fixa, então um trim indevido mascararia erros reais
// de tamanho de linha.
func splitLinhas(conteudo []byte) []string {
	texto := strings.ReplaceAll(string(conteudo), "\r\n", "\n")
	texto = strings.TrimSuffix(texto, "\n")

	if texto == "" {
		return nil
	}

	return strings.Split(texto, "\n")
}

// detectarTipo identifica CNAB 240 ou CNAB 400 pelo tamanho de linha mais
// frequente no arquivo — usar a moda (e não só a primeira linha) evita que
// uma única linha corrompida atrapalhe a detecção do layout.
func detectarTipo(linhas []string) (TipoCnab, error) {
	contagem := map[int]int{}
	for _, l := range linhas {
		contagem[len(l)]++
	}

	tamanhoDominante, maiorContagem := 0, 0
	for tamanho, qtd := range contagem {
		if qtd > maiorContagem {
			tamanhoDominante, maiorContagem = tamanho, qtd
		}
	}

	switch tamanhoDominante {
	case 240:
		return CNAB240, nil
	case 400:
		return CNAB400, nil
	default:
		return "", fmt.Errorf("não foi possível identificar o layout do arquivo: tamanho de linha predominante é %d caracteres (esperado 240 para CNAB 240 ou 400 para CNAB 400)", tamanhoDominante)
	}
}

func parseNumero(campo string) (int, error) {
	campo = strings.TrimSpace(campo)
	if campo == "" {
		return 0, errors.New("Numeric field is empty")
	}
	return strconv.Atoi(campo)
}

// campo extrai com segurança um trecho de largura fixa da linha (posições
// 1-indexadas, como na documentação FEBRABAN), devolvendo string vazia se a
// linha for curta demais para conter aquele campo.
func campo(linha string, inicio, fim int) string {
	if inicio < 1 {
		inicio = 1
	}
	if inicio > len(linha) {
		return ""
	}
	if fim > len(linha) {
		fim = len(linha)
	}
	return linha[inicio-1 : fim]
}

// ValidarRemessa lê um arquivo de remessa CNAB 240 ou CNAB 400, detecta o
// layout automaticamente e retorna a lista de erros estruturais encontrados
// — cada um já apontando a linha exata do arquivo original.
func ValidarRemessa(conteudo []byte) (*Resultado, error) {
	linhas := splitLinhas(conteudo)
	if len(linhas) == 0 {
		return nil, errors.New("Empty file")
	}

	tipo, err := detectarTipo(linhas)
	if err != nil {
		return nil, err
	}

	resultado := &Resultado{
		TipoCnab:    tipo,
		TotalLinhas: len(linhas),
		Erros:       []Erro{},
	}

	if tipo == CNAB400 {
		resultado.Erros = validarCNAB400(linhas)
	} else {
		erros, quantidadeLotes := validarCNAB240(linhas)
		resultado.Erros = erros
		resultado.QuantidadeLotes = quantidadeLotes
	}

	resultado.Valido = len(resultado.Erros) == 0

	return resultado, nil
}

// validarCNAB400 confere o tamanho de cada linha e o número sequencial do
// registro (posições 395-400), que numera TODAS as linhas do arquivo em
// sequência — do header (1) ao trailer (última). Uma quebra nessa sequência
// é o sinal mais confiável de linha faltante, duplicada ou fora de ordem.
func validarCNAB400(linhas []string) []Erro {
	var erros []Erro

	primeiroTipo := campo(linhas[0], 1, 1)
	if primeiroTipo != "0" {
		erros = append(erros, Erro{
			Linha:     1,
			Categoria: "estrutura",
			Descricao: fmt.Sprintf("a primeira linha deveria ser o Header de Arquivo (tipo 0), mas o tipo de registro encontrado foi %q", primeiroTipo),
		})
	}

	ultimaLinha := len(linhas)
	ultimoTipo := campo(linhas[ultimaLinha-1], 1, 1)
	if ultimoTipo != "9" {
		erros = append(erros, Erro{
			Linha:     ultimaLinha,
			Categoria: "estrutura",
			Descricao: fmt.Sprintf("a última linha deveria ser o Trailer de Arquivo (tipo 9), mas o tipo de registro encontrado foi %q — o arquivo pode estar truncado", ultimoTipo),
		})
	}

	esperado := 0
	for i, linha := range linhas {
		numLinha := i + 1

		if len(linha) != 400 {
			erros = append(erros, Erro{
				Linha:     numLinha,
				Categoria: "tamanho",
				Descricao: fmt.Sprintf("linha com %d caracteres, esperado 400 (CNAB 400)", len(linha)),
			})
			continue
		}

		esperado++

		seq, err := parseNumero(campo(linha, 395, 400))
		if err != nil {
			erros = append(erros, Erro{
				Linha:     numLinha,
				Categoria: "sequencia",
				Descricao: fmt.Sprintf("não foi possível ler o número sequencial do registro (posições 395-400): %q", campo(linha, 395, 400)),
			})
			continue
		}

		if seq != esperado {
			erros = append(erros, sequenciaQuebradaErro(numLinha, esperado, seq))
			esperado = seq // ressincroniza para não gerar uma cascata de erros repetidos
		}
	}

	return erros
}

// validarCNAB240 confere tamanho de linha, a ordem dos tipos de registro
// (Header de Arquivo → [Header de Lote → Detalhes → Trailer de Lote]* →
// Trailer de Arquivo) e o número sequencial dos registros de detalhe dentro
// de cada lote (posições 9-13, reiniciado a cada novo lote). Também
// confronta a quantidade de registros de cada Trailer de Lote e a
// quantidade total declarada no Trailer de Arquivo contra o que foi
// efetivamente contado no arquivo.
func validarCNAB240(linhas []string) ([]Erro, int) {
	var erros []Erro

	const (
		estadoInicial        = iota // antes de qualquer registro
		estadoAposHeaderArq         // aguardando Header de Lote ou Trailer de Arquivo
		estadoDentroDoLote          // aguardando Detalhe ou Trailer de Lote
		estadoAposTrailerArq        // não deveria haver mais linhas
	)

	estado := estadoInicial
	quantidadeLotes := 0
	seqDetalheAtual := 0
	inicioLoteLinha := 0
	trailerArquivoEncontrado := false

	for i, linha := range linhas {
		numLinha := i + 1

		if len(linha) != 240 {
			erros = append(erros, Erro{
				Linha:     numLinha,
				Categoria: "tamanho",
				Descricao: fmt.Sprintf("linha com %d caracteres, esperado 240 (CNAB 240)", len(linha)),
			})
			continue
		}

		tipo := campo(linha, 8, 8)

		switch estado {
		case estadoInicial:
			if tipo != "0" {
				erros = append(erros, Erro{
					Linha:     numLinha,
					Categoria: "estrutura",
					Descricao: fmt.Sprintf("a primeira linha deveria ser o Header de Arquivo (tipo 0), mas o tipo de registro encontrado foi %q", tipo),
				})
			}
			estado = estadoAposHeaderArq

		case estadoAposHeaderArq:
			switch tipo {
			case "1":
				quantidadeLotes++
				seqDetalheAtual = 0
				inicioLoteLinha = numLinha
				estado = estadoDentroDoLote
			case "9":
				validarTrailerArquivo(linha, numLinha, len(linhas), quantidadeLotes, &erros)
				trailerArquivoEncontrado = true
				estado = estadoAposTrailerArq
			default:
				erros = append(erros, Erro{
					Linha:     numLinha,
					Categoria: "estrutura",
					Descricao: fmt.Sprintf("esperado Header de Lote (tipo 1) ou Trailer de Arquivo (tipo 9), mas encontrado tipo %q — pode faltar o Header de Lote", tipo),
				})
				// tenta se recuperar tratando a linha como início de um lote
				quantidadeLotes++
				seqDetalheAtual = 0
				inicioLoteLinha = numLinha
				estado = estadoDentroDoLote
			}

		case estadoDentroDoLote:
			switch tipo {
			case "3":
				seqDetalheAtual++
				seq, err := parseNumero(campo(linha, 9, 13))
				if err != nil {
					erros = append(erros, Erro{
						Linha:     numLinha,
						Categoria: "sequencia",
						Descricao: fmt.Sprintf("não foi possível ler o número sequencial do registro no lote (posições 9-13): %q", campo(linha, 9, 13)),
					})
					continue
				}
				if seq != seqDetalheAtual {
					erros = append(erros, sequenciaQuebradaErro(numLinha, seqDetalheAtual, seq))
					seqDetalheAtual = seq // ressincroniza
				}
			case "5":
				declarado, err := parseNumero(campo(linha, 18, 23))
				esperadoNoLote := seqDetalheAtual + 2 // header de lote + detalhes + este trailer
				if err != nil {
					erros = append(erros, Erro{
						Linha:     numLinha,
						Categoria: "sequencia",
						Descricao: fmt.Sprintf("não foi possível ler a quantidade de registros do lote (posições 18-23): %q", campo(linha, 18, 23)),
					})
				} else if declarado != esperadoNoLote {
					erros = append(erros, Erro{
						Linha:     numLinha,
						Categoria: "sequencia",
						Descricao: fmt.Sprintf("Trailer de Lote declara %d registros no lote, mas foram contados %d (lote iniciado na linha %d) — provável linha faltante ou duplicada dentro do lote", declarado, esperadoNoLote, inicioLoteLinha),
					})
				}
				estado = estadoAposHeaderArq
			default:
				erros = append(erros, Erro{
					Linha:     numLinha,
					Categoria: "estrutura",
					Descricao: fmt.Sprintf("esperado Detalhe (tipo 3) ou Trailer de Lote (tipo 5), mas encontrado tipo %q — pode faltar o Trailer de Lote", tipo),
				})
			}

		case estadoAposTrailerArq:
			erros = append(erros, Erro{
				Linha:     numLinha,
				Categoria: "estrutura",
				Descricao: "linha encontrada após o Trailer de Arquivo — o arquivo deveria terminar ali",
			})
		}
	}

	if !trailerArquivoEncontrado {
		erros = append(erros, Erro{
			Linha:     len(linhas),
			Categoria: "estrutura",
			Descricao: "o arquivo não termina com um Trailer de Arquivo (tipo 9) — o arquivo pode estar truncado ou com a última linha faltando",
		})
	}

	return erros, quantidadeLotes
}

func validarTrailerArquivo(linha string, numLinha, totalLinhas, quantidadeLotesContada int, erros *[]Erro) {
	lotesDeclarados, err := parseNumero(campo(linha, 18, 23))
	if err != nil {
		*erros = append(*erros, Erro{
			Linha:     numLinha,
			Categoria: "sequencia",
			Descricao: fmt.Sprintf("não foi possível ler a quantidade de lotes do arquivo (posições 18-23): %q", campo(linha, 18, 23)),
		})
	} else if lotesDeclarados != quantidadeLotesContada {
		*erros = append(*erros, Erro{
			Linha:     numLinha,
			Categoria: "sequencia",
			Descricao: fmt.Sprintf("Trailer de Arquivo declara %d lote(s), mas foram contados %d — provável lote inteiro faltando", lotesDeclarados, quantidadeLotesContada),
		})
	}

	registrosDeclarados, err := parseNumero(campo(linha, 24, 29))
	if err != nil {
		*erros = append(*erros, Erro{
			Linha:     numLinha,
			Categoria: "sequencia",
			Descricao: fmt.Sprintf("não foi possível ler a quantidade de registros do arquivo (posições 24-29): %q", campo(linha, 24, 29)),
		})
	} else if registrosDeclarados != totalLinhas {
		*erros = append(*erros, Erro{
			Linha:     numLinha,
			Categoria: "sequencia",
			Descricao: fmt.Sprintf("Trailer de Arquivo declara %d registro(s) no total, mas o arquivo tem %d linha(s) — provável linha faltante", registrosDeclarados, totalLinhas),
		})
	}
}

func sequenciaQuebradaErro(numLinha, esperado, encontrado int) Erro {
	if encontrado > esperado {
		faltando := encontrado - esperado
		return Erro{
			Linha:     numLinha,
			Categoria: "sequencia",
			Descricao: fmt.Sprintf("quebra de sequência: número esperado era %d, mas foi encontrado %d — parece(m) estar faltando %d linha(s) antes desta", esperado, encontrado, faltando),
		}
	}

	return Erro{
		Linha:     numLinha,
		Categoria: "sequencia",
		Descricao: fmt.Sprintf("sequência fora de ordem: número esperado era %d, mas foi encontrado %d (menor que o esperado — possível linha duplicada ou fora de ordem)", esperado, encontrado),
	}
}
